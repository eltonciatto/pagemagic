import { PromptTemplate } from '../types/prompt';
export declare class TemplateService {
    private templates;
    constructor();
    getTemplates(filters?: {
        category?: string;
        tags?: string[];
    }): Promise<PromptTemplate[]>;
    getTemplate(templateId: string): Promise<PromptTemplate | null>;
    applyTemplate(templateId: string, variables: Record<string, any>): Promise<string>;
    createTemplate(template: Omit<PromptTemplate, 'id' | 'createdAt' | 'updatedAt'>): Promise<PromptTemplate>;
    private initializeDefaultTemplates;
}
//# sourceMappingURL=TemplateService.d.ts.map