export declare const config: {
    port: number;
    environment: string;
    openai: {
        apiKey: string;
        model: string;
        temperature: number;
    };
    anthropic: {
        apiKey: string;
        model: string;
        temperature: number;
    };
    vllm: {
        baseUrl: string;
        model: string;
        temperature: number;
    };
    redis: {
        host: string;
        port: number;
        password: string;
        db: number;
    };
    nats: {
        url: string;
        user: string;
        password: string;
    };
    cors: {
        origins: string[];
    };
    rateLimit: {
        windowMs: number;
        max: number;
    };
    jwt: {
        secret: string;
    };
    logging: {
        level: string;
        format: string;
    };
    metrics: {
        enabled: boolean;
        prefix: string;
    };
};
//# sourceMappingURL=index.d.ts.map