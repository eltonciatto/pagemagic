export interface SiteGenerationResponse {
    site: {
        title: string;
        description: string;
        sections: Section[];
        theme: Theme;
        metadata: Metadata;
    };
}
export interface Section {
    id: string;
    type: 'hero' | 'features' | 'testimonials' | 'cta' | 'about' | 'contact';
    title: string;
    content: string;
    cta?: CallToAction;
    features?: Feature[];
    testimonials?: Testimonial[];
    order: number;
}
export interface Feature {
    title: string;
    description: string;
    icon?: string;
}
export interface Testimonial {
    quote: string;
    author: string;
    role?: string;
    company?: string;
}
export interface CallToAction {
    text: string;
    link: string;
    style: 'primary' | 'secondary' | 'outline';
}
export interface Theme {
    primaryColor: string;
    secondaryColor: string;
    fontFamily: 'sans' | 'serif' | 'mono';
    layout: 'modern' | 'classic' | 'minimal';
}
export interface Metadata {
    industry?: string;
    target_audience?: string;
    tone?: 'professional' | 'casual' | 'friendly' | 'authoritative';
    language: string;
}
export interface SiteGenerationRequest {
    description: string;
    userId: string;
    projectId?: string;
    preferences?: {
        industry?: string;
        style?: string;
        features?: string[];
    };
}
export interface GenerationApiResponse {
    success: boolean;
    data?: SiteGenerationResponse;
    error?: string;
    metadata: {
        requestId: string;
        tokensUsed: number;
        processingTime: number;
        model: string;
    };
}
//# sourceMappingURL=generation.d.ts.map