// Simple MVP Builder Service in Rust
use axum::{
    extract::Json,
    http::StatusCode,
    response::Json as ResponseJson,
    routing::{get, post},
    Router,
};
use serde::{Deserialize, Serialize};
use tower_http::cors::CorsLayer;

// Site structure from prompt-svc
#[derive(Debug, Deserialize, Serialize)]
struct SiteData {
    site: Site,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
struct Site {
    title: String,
    description: String,
    sections: Vec<Section>,
    theme: Theme,
    metadata: Metadata,
}

#[derive(Debug, Deserialize, Serialize)]
struct Section {
    id: String,
    #[serde(rename = "type")]
    section_type: String,
    title: String,
    content: String,
    order: i32,
    cta: Option<CallToAction>,
    features: Option<Vec<Feature>>,
    testimonials: Option<Vec<Testimonial>>,
}

#[derive(Debug, Deserialize, Serialize)]
struct CallToAction {
    text: String,
    link: String,
    style: String,
}

#[derive(Debug, Deserialize, Serialize)]
struct Feature {
    title: String,
    description: String,
    icon: Option<String>,
}

#[derive(Debug, Deserialize, Serialize)]
struct Testimonial {
    quote: String,
    author: String,
    role: Option<String>,
    company: Option<String>,
}

#[derive(Debug, Deserialize, Serialize)]
struct Theme {
    #[serde(rename = "primaryColor")]
    primary_color: String,
    #[serde(rename = "secondaryColor")]
    secondary_color: String,
    #[serde(rename = "fontFamily")]
    font_family: String,
    layout: String,
}

#[derive(Debug, Deserialize, Serialize)]
struct Metadata {
    industry: Option<String>,
    target_audience: Option<String>,
    tone: Option<String>,
    language: String,
}

// Request/Response types
#[derive(Debug, Deserialize)]
struct BuildRequest {
    data: SiteData,
}

#[derive(Debug, Serialize)]
struct BuildResponse {
    success: bool,
    data: Option<BuildData>,
    error: Option<String>,
    metadata: BuildMetadata,
}

#[derive(Debug, Serialize)]
struct BuildData {
    html: String,
    css: String,
    #[serde(rename = "buildId")]
    build_id: String,
}

#[derive(Debug, Serialize)]
struct BuildMetadata {
    #[serde(rename = "requestId")]
    request_id: String,
    #[serde(rename = "processingTime")]
    processing_time: i64,
    #[serde(rename = "sectionsProcessed")]
    sections_processed: usize,
}

// HTML template builder
struct HtmlBuilder {
    site: Site,
}

impl HtmlBuilder {
    fn new(site: Site) -> Self {
        Self { site }
    }

    fn build(&self) -> (String, String) {
        let html = self.generate_html();
        let css = self.generate_css();
        (html, css)
    }

    fn generate_html(&self) -> String {
        let sections_html = self.site.sections
            .iter()
            .map(|section| self.build_section(section))
            .collect::<Vec<_>>()
            .join("\n");

        format!(
            r#"<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{}</title>
    <meta name="description" content="{}">
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {{
            theme: {{
                extend: {{
                    colors: {{
                        primary: '{}',
                        secondary: '{}'
                    }}
                }}
            }}
        }}
    </script>
</head>
<body class="font-{} bg-white">
    {}
</body>
</html>"#,
            self.site.title,
            self.site.description,
            self.site.theme.primary_color,
            self.site.theme.secondary_color,
            self.site.theme.font_family,
            sections_html
        )
    }

    fn build_section(&self, section: &Section) -> String {
        match section.section_type.as_str() {
            "hero" => self.build_hero_section(section),
            "features" => self.build_features_section(section),
            "testimonials" => self.build_testimonials_section(section),
            "cta" => self.build_cta_section(section),
            _ => self.build_generic_section(section),
        }
    }

    fn build_hero_section(&self, section: &Section) -> String {
        let cta_html = if let Some(cta) = &section.cta {
            format!(
                r#"<a href="{}" class="inline-block bg-primary text-white px-8 py-3 rounded-lg font-semibold hover:opacity-90 transition-opacity">{}</a>"#,
                cta.link, cta.text
            )
        } else {
            String::new()
        };

        format!(
            r#"<section class="bg-gradient-to-br from-gray-50 to-gray-100 py-20">
    <div class="container mx-auto px-6 text-center">
        <h1 class="text-5xl font-bold text-gray-900 mb-6">{}</h1>
        <p class="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">{}</p>
        {}
    </div>
</section>"#,
            section.title, section.content, cta_html
        )
    }

    fn build_features_section(&self, section: &Section) -> String {
        let features_html = if let Some(features) = &section.features {
            features
                .iter()
                .map(|feature| {
                    format!(
                        r#"<div class="bg-white p-6 rounded-lg shadow-md">
                            <h3 class="text-xl font-semibold text-gray-900 mb-3">{}</h3>
                            <p class="text-gray-600">{}</p>
                        </div>"#,
                        feature.title, feature.description
                    )
                })
                .collect::<Vec<_>>()
                .join("\n")
        } else {
            String::new()
        };

        format!(
            r#"<section class="py-16 bg-white">
    <div class="container mx-auto px-6">
        <div class="text-center mb-12">
            <h2 class="text-4xl font-bold text-gray-900 mb-4">{}</h2>
            <p class="text-xl text-gray-600">{}</p>
        </div>
        <div class="grid md:grid-cols-3 gap-8">
            {}
        </div>
    </div>
</section>"#,
            section.title, section.content, features_html
        )
    }

    fn build_testimonials_section(&self, section: &Section) -> String {
        let testimonials_html = if let Some(testimonials) = &section.testimonials {
            testimonials
                .iter()
                .map(|testimonial| {
                    let role_company = match (&testimonial.role, &testimonial.company) {
                        (Some(role), Some(company)) => format!("{} at {}", role, company),
                        (Some(role), None) => role.clone(),
                        (None, Some(company)) => company.clone(),
                        (None, None) => String::new(),
                    };

                    format!(
                        r#"<div class="bg-white p-6 rounded-lg shadow-md">
                            <p class="text-gray-700 mb-4 italic">"{}"</p>
                            <div class="font-semibold text-gray-900">{}</div>
                            <div class="text-gray-600 text-sm">{}</div>
                        </div>"#,
                        testimonial.quote, testimonial.author, role_company
                    )
                })
                .collect::<Vec<_>>()
                .join("\n")
        } else {
            String::new()
        };

        format!(
            r#"<section class="py-16 bg-gray-50">
    <div class="container mx-auto px-6">
        <div class="text-center mb-12">
            <h2 class="text-4xl font-bold text-gray-900 mb-4">{}</h2>
            <p class="text-xl text-gray-600">{}</p>
        </div>
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            {}
        </div>
    </div>
</section>"#,
            section.title, section.content, testimonials_html
        )
    }

    fn build_cta_section(&self, section: &Section) -> String {
        let cta_html = if let Some(cta) = &section.cta {
            format!(
                r#"<a href="{}" class="inline-block bg-primary text-white px-8 py-4 rounded-lg font-semibold text-lg hover:opacity-90 transition-opacity">{}</a>"#,
                cta.link, cta.text
            )
        } else {
            String::new()
        };

        format!(
            r#"<section class="py-20 bg-primary text-white">
    <div class="container mx-auto px-6 text-center">
        <h2 class="text-4xl font-bold mb-6">{}</h2>
        <p class="text-xl mb-8 opacity-90">{}</p>
        {}
    </div>
</section>"#,
            section.title, section.content, cta_html
        )
    }

    fn build_generic_section(&self, section: &Section) -> String {
        format!(
            r#"<section class="py-16">
    <div class="container mx-auto px-6">
        <h2 class="text-3xl font-bold text-gray-900 mb-6">{}</h2>
        <div class="prose max-w-none">{}</div>
    </div>
</section>"#,
            section.title, section.content
        )
    }

    fn generate_css(&self) -> String {
        // Additional custom CSS if needed
        format!(
            r#"/* Custom styles for {} */
.container {{
    max-width: 1200px;
}}

.prose {{
    line-height: 1.7;
}}

/* Custom theme adjustments */
:root {{
    --primary-color: {};
    --secondary-color: {};
}}

/* Responsive adjustments */
@media (max-width: 768px) {{
    .text-5xl {{
        font-size: 2.5rem;
    }}
    
    .text-4xl {{
        font-size: 2rem;
    }}
}}"#,
            self.site.title, self.site.theme.primary_color, self.site.theme.secondary_color
        )
    }
}

// API Handlers
async fn health() -> ResponseJson<serde_json::Value> {
    ResponseJson(serde_json::json!({
        "status": "healthy",
        "service": "builder-svc",
        "timestamp": chrono::Utc::now().to_rfc3339(),
        "version": "mvp-0.1.0"
    }))
}

async fn build_site(Json(request): Json<BuildRequest>) -> Result<ResponseJson<BuildResponse>, StatusCode> {
    let start_time = std::time::Instant::now();
    let request_id = format!("build_{}", chrono::Utc::now().timestamp_millis());
    
    println!("Building site: {}", request.data.site.title);

    // Build HTML/CSS
    let builder = HtmlBuilder::new(request.data.site.clone());
    let (html, css) = builder.build();
    
    let build_id = format!("site_{}", uuid::Uuid::new_v4());
    let processing_time = start_time.elapsed().as_millis() as i64;

    let response = BuildResponse {
        success: true,
        data: Some(BuildData {
            html,
            css,
            build_id,
        }),
        error: None,
        metadata: BuildMetadata {
            request_id,
            processing_time,
            sections_processed: request.data.site.sections.len(),
        },
    };

    Ok(ResponseJson(response))
}

#[tokio::main]
async fn main() {
    // Initialize tracing
    tracing_subscriber::fmt::init();

    // Build our application with routes
    let app = Router::new()
        .route("/v1/health", get(health))
        .route("/v1/build", post(build_site))
        .layer(CorsLayer::permissive());

    // Run it with hyper on port 3002
    let port = std::env::var("PORT").unwrap_or_else(|_| "3002".to_string());
    let addr = format!("0.0.0.0:{}", port);
    
    println!("🚀 Builder service (MVP) running on port {}", port);
    println!("📋 Test endpoint: http://localhost:{}/v1/build", port);
    println!("❤️  Health check: http://localhost:{}/v1/health", port);

    let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
