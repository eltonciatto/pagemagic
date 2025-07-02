import httpx
import logging
from typing import List, Optional
from ..models.domain import DomainAvailability, DomainRegistration, Domain
from ..config.settings import Settings

logger = logging.getLogger(__name__)

class NamecheapService:
    def __init__(self, settings: Settings):
        self.api_key = settings.namecheap_api_key
        self.username = settings.namecheap_username
        self.api_user = settings.namecheap_api_user
        self.sandbox = settings.namecheap_sandbox
        
        self.base_url = "https://api.sandbox.namecheap.com/xml.response" if self.sandbox else "https://api.namecheap.com/xml.response"
    
    async def check_availability(self, domains: List[str]) -> List[DomainAvailability]:
        """Check domain availability using Namecheap API"""
        try:
            domain_list = ",".join(domains)
            params = {
                "ApiUser": self.api_user,
                "ApiKey": self.api_key,
                "UserName": self.username,
                "Command": "namecheap.domains.check",
                "ClientIp": "127.0.0.1",  # In production, get actual client IP
                "DomainList": domain_list
            }
            
            async with httpx.AsyncClient() as client:
                response = await client.get(self.base_url, params=params)
                response.raise_for_status()
                
                # Parse XML response and convert to our models
                # This is simplified - real implementation would parse XML
                results = []
                for domain in domains:
                    results.append(DomainAvailability(
                        domain=domain,
                        available=True,  # Placeholder
                        price=10.99,
                        currency="USD"
                    ))
                
                return results
                
        except Exception as e:
            logger.error(f"Error checking domain availability: {e}")
            raise
    
    async def register_domain(self, registration: DomainRegistration) -> Domain:
        """Register a domain through Namecheap"""
        try:
            params = {
                "ApiUser": self.api_user,
                "ApiKey": self.api_key,
                "UserName": self.username,
                "Command": "namecheap.domains.create",
                "ClientIp": "127.0.0.1",
                "DomainName": registration.domain,
                "Years": registration.years,
                # Add registrant info params...
            }
            
            async with httpx.AsyncClient() as client:
                response = await client.post(self.base_url, data=params)
                response.raise_for_status()
                
                # Parse response and create Domain object
                # This is simplified
                from datetime import datetime, timedelta
                domain = Domain(
                    id=f"dom_{registration.domain.replace('.', '_')}",
                    name=registration.domain,
                    status="registered",
                    user_id=registration.user_id,
                    project_id=registration.project_id,
                    registrar="namecheap",
                    expires_at=datetime.utcnow() + timedelta(days=365 * registration.years),
                    auto_renew=registration.auto_renew,
                    nameservers=[],
                    created_at=datetime.utcnow(),
                    updated_at=datetime.utcnow()
                )
                
                return domain
                
        except Exception as e:
            logger.error(f"Error registering domain {registration.domain}: {e}")
            raise
