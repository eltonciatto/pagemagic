import { PromptAnalytics } from '../types/prompt';
interface UserAnalytics {
    totalRequests: number;
    successfulRequests: number;
    failedRequests: number;
    totalTokensUsed: number;
    totalCost: number;
    averageProcessingTime: number;
    topModelsUsed: {
        model: string;
        count: number;
    }[];
    requestsByDay: {
        date: string;
        count: number;
    }[];
    errorTypes: {
        type: string;
        count: number;
    }[];
}
export declare class AnalyticsService {
    private analytics;
    trackPromptRequest(data: PromptAnalytics): Promise<void>;
    getUserAnalytics(userId: string, options?: {
        startDate?: Date;
        endDate?: Date;
    }): Promise<UserAnalytics>;
    getGlobalAnalytics(options?: {
        startDate?: Date;
        endDate?: Date;
    }): Promise<{
        totalUsers: number;
        totalRequests: number;
        totalTokensUsed: number;
        totalCost: number;
        successRate: number;
        topUsers: {
            userId: string;
            requestCount: number;
        }[];
        topModels: {
            model: string;
            count: number;
        }[];
    }>;
    getUserUsage(userId: string, timeframe?: 'day' | 'week' | 'month'): Promise<{
        requestCount: number;
        tokenCount: number;
        cost: number;
        limit?: number;
        isOverLimit: boolean;
    }>;
    exportUserAnalytics(userId: string, format?: 'json' | 'csv'): Promise<string>;
    clearUserAnalytics(userId: string): Promise<boolean>;
}
export {};
//# sourceMappingURL=AnalyticsService.d.ts.map