import os
import logging
from typing import Optional
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    # App
    app_name: str = "domain-svc"
    version: str = "1.0.0"
    debug: bool = False
    port: int = 8080
    
    # Database
    database_url: str = "postgresql://user:pass@localhost:5432/pagemagic"
    
    # Domain providers
    namecheap_api_key: Optional[str] = None
    namecheap_username: Optional[str] = None
    namecheap_api_user: Optional[str] = None
    namecheap_sandbox: bool = True
    
    # DNS
    cloudflare_api_token: Optional[str] = None
    cloudflare_zone_id: Optional[str] = None
    
    # ACME
    acme_email: str = "admin@pagemagic.com"
    acme_directory_url: str = "https://acme-v02.api.letsencrypt.org/directory"
    acme_staging: bool = True
    
    # NATS
    nats_url: str = "nats://localhost:4222"
    
    # Monitoring
    metrics_enabled: bool = True
    log_level: str = "INFO"
    
    class Config:
        env_file = ".env"
        case_sensitive = False

def get_settings() -> Settings:
    return Settings()

def setup_logging(settings: Settings):
    log_level = getattr(logging, settings.log_level.upper())
    logging.basicConfig(
        level=log_level,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    )
