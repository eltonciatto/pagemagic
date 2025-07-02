import { PromptAnalytics } from '../types/prompt';
import { logger } from '../utils/logger';

interface UserAnalytics {
  totalRequests: number;
  successfulRequests: number;
  failedRequests: number;
  totalTokensUsed: number;
  totalCost: number;
  averageProcessingTime: number;
  topModelsUsed: { model: string; count: number }[];
  requestsByDay: { date: string; count: number }[];
  errorTypes: { type: string; count: number }[];
}

export class AnalyticsService {
  private analytics: Map<string, PromptAnalytics[]> = new Map();

  async trackPromptRequest(data: PromptAnalytics): Promise<void> {
    try {
      const userAnalytics = this.analytics.get(data.userId) || [];
      userAnalytics.push(data);
      this.analytics.set(data.userId, userAnalytics);

      // Log for monitoring
      logger.info(`Tracked prompt request: user=${data.userId}, tokens=${data.tokensUsed}, cost=${data.cost}, success=${data.success}`);
    } catch (error) {
      logger.error(`Failed to track analytics: ${error}`);
    }
  }

  async getUserAnalytics(
    userId: string,
    options: { startDate?: Date; endDate?: Date } = {}
  ): Promise<UserAnalytics> {
    const userRecords = this.analytics.get(userId) || [];
    
    // Filter by date range if provided
    const filteredRecords = userRecords.filter(record => {
      if (options.startDate && record.timestamp < options.startDate) return false;
      if (options.endDate && record.timestamp > options.endDate) return false;
      return true;
    });

    const totalRequests = filteredRecords.length;
    const successfulRequests = filteredRecords.filter(r => r.success).length;
    const failedRequests = totalRequests - successfulRequests;
    const totalTokensUsed = filteredRecords.reduce((sum, r) => sum + r.tokensUsed, 0);
    const totalCost = filteredRecords.reduce((sum, r) => sum + r.cost, 0);
    
    const averageProcessingTime = totalRequests > 0 
      ? filteredRecords.reduce((sum, r) => sum + r.processingTime, 0) / totalRequests 
      : 0;

    // Top models used
    const modelCounts: Record<string, number> = {};
    filteredRecords.forEach(record => {
      modelCounts[record.model] = (modelCounts[record.model] || 0) + 1;
    });
    
    const topModelsUsed = Object.entries(modelCounts)
      .map(([model, count]) => ({ model, count }))
      .sort((a, b) => b.count - a.count)
      .slice(0, 5);

    // Requests by day
    const dailyCounts: Record<string, number> = {};
    filteredRecords.forEach(record => {
      const date = record.timestamp.toISOString().split('T')[0];
      dailyCounts[date] = (dailyCounts[date] || 0) + 1;
    });
    
    const requestsByDay = Object.entries(dailyCounts)
      .map(([date, count]) => ({ date, count }))
      .sort((a, b) => a.date.localeCompare(b.date));

    // Error types
    const errorCounts: Record<string, number> = {};
    filteredRecords
      .filter(r => !r.success && r.errorType)
      .forEach(record => {
        errorCounts[record.errorType!] = (errorCounts[record.errorType!] || 0) + 1;
      });
    
    const errorTypes = Object.entries(errorCounts)
      .map(([type, count]) => ({ type, count }))
      .sort((a, b) => b.count - a.count);

    return {
      totalRequests,
      successfulRequests,
      failedRequests,
      totalTokensUsed,
      totalCost,
      averageProcessingTime,
      topModelsUsed,
      requestsByDay,
      errorTypes,
    };
  }

  async getGlobalAnalytics(options: { startDate?: Date; endDate?: Date } = {}): Promise<{
    totalUsers: number;
    totalRequests: number;
    totalTokensUsed: number;
    totalCost: number;
    successRate: number;
    topUsers: { userId: string; requestCount: number }[];
    topModels: { model: string; count: number }[];
  }> {
    const allRecords: PromptAnalytics[] = [];
    
    // Collect all records from all users
    for (const userRecords of this.analytics.values()) {
      const filteredRecords = userRecords.filter(record => {
        if (options.startDate && record.timestamp < options.startDate) return false;
        if (options.endDate && record.timestamp > options.endDate) return false;
        return true;
      });
      allRecords.push(...filteredRecords);
    }

    const totalUsers = this.analytics.size;
    const totalRequests = allRecords.length;
    const totalTokensUsed = allRecords.reduce((sum, r) => sum + r.tokensUsed, 0);
    const totalCost = allRecords.reduce((sum, r) => sum + r.cost, 0);
    const successfulRequests = allRecords.filter(r => r.success).length;
    const successRate = totalRequests > 0 ? successfulRequests / totalRequests : 0;

    // Top users by request count
    const userCounts: Record<string, number> = {};
    allRecords.forEach(record => {
      userCounts[record.userId] = (userCounts[record.userId] || 0) + 1;
    });
    
    const topUsers = Object.entries(userCounts)
      .map(([userId, requestCount]) => ({ userId, requestCount }))
      .sort((a, b) => b.requestCount - a.requestCount)
      .slice(0, 10);

    // Top models
    const modelCounts: Record<string, number> = {};
    allRecords.forEach(record => {
      modelCounts[record.model] = (modelCounts[record.model] || 0) + 1;
    });
    
    const topModels = Object.entries(modelCounts)
      .map(([model, count]) => ({ model, count }))
      .sort((a, b) => b.count - a.count)
      .slice(0, 10);

    return {
      totalUsers,
      totalRequests,
      totalTokensUsed,
      totalCost,
      successRate,
      topUsers,
      topModels,
    };
  }

  async getUserUsage(userId: string, timeframe: 'day' | 'week' | 'month' = 'month'): Promise<{
    requestCount: number;
    tokenCount: number;
    cost: number;
    limit?: number;
    isOverLimit: boolean;
  }> {
    const now = new Date();
    let startDate: Date;

    switch (timeframe) {
      case 'day':
        startDate = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        break;
      case 'week':
        startDate = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
        break;
      case 'month':
        startDate = new Date(now.getFullYear(), now.getMonth(), 1);
        break;
    }

    const userRecords = this.analytics.get(userId) || [];
    const periodRecords = userRecords.filter(record => record.timestamp >= startDate);

    const requestCount = periodRecords.length;
    const tokenCount = periodRecords.reduce((sum, r) => sum + r.tokensUsed, 0);
    const cost = periodRecords.reduce((sum, r) => sum + r.cost, 0);

    // Example limits - these should come from user subscription/plan
    const limits = {
      day: { requests: 100, tokens: 50000 },
      week: { requests: 500, tokens: 250000 },
      month: { requests: 2000, tokens: 1000000 },
    };

    const limit = limits[timeframe];
    const isOverLimit = requestCount >= limit.requests || tokenCount >= limit.tokens;

    return {
      requestCount,
      tokenCount,
      cost,
      limit: limit.requests,
      isOverLimit,
    };
  }

  async exportUserAnalytics(userId: string, format: 'json' | 'csv' = 'json'): Promise<string> {
    const userRecords = this.analytics.get(userId) || [];
    
    if (format === 'csv') {
      const headers = 'timestamp,prompt,tokensUsed,cost,processingTime,success,model,errorType\n';
      const rows = userRecords.map(record => 
        `${record.timestamp.toISOString()},${JSON.stringify(record.prompt)},${record.tokensUsed},${record.cost},${record.processingTime},${record.success},${record.model},${record.errorType || ''}`
      ).join('\n');
      
      return headers + rows;
    }

    return JSON.stringify(userRecords, null, 2);
  }

  async clearUserAnalytics(userId: string): Promise<boolean> {
    try {
      this.analytics.delete(userId);
      return true;
    } catch (error) {
      logger.error(`Failed to clear user analytics: ${error}`);
      return false;
    }
  }
}
