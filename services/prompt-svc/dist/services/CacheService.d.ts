export declare class CacheService {
    private redis;
    private isConnected;
    constructor();
    get<T>(key: string): Promise<T | null>;
    set<T>(key: string, value: T, ttlSeconds?: number): Promise<boolean>;
    delete(key: string): Promise<boolean>;
    exists(key: string): Promise<boolean>;
    increment(key: string, ttlSeconds?: number): Promise<number>;
    flush(): Promise<boolean>;
    disconnect(): Promise<void>;
    isHealthy(): boolean;
}
//# sourceMappingURL=CacheService.d.ts.map