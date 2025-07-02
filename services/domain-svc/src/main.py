from fastapi import FastAPI, HTTPException, Depends
from fastapi.middleware.cors import CORSMiddleware
from typing import List
import logging

from .config.settings import Settings, get_settings, setup_logging
from .models.domain import DomainAvailability, DomainRegistration, Domain, DNSRecord
from .services.namecheap import NamecheapService
from .services.dns import CloudflareService
from .services.acme import ACMEService

# Setup logging
settings = get_settings()
setup_logging(settings)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Domain Service",
    description="Page Magic Domain Management Service",
    version="1.0.0"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure appropriately for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Initialize services
namecheap_service = NamecheapService(settings)
dns_service = CloudflareService(settings)
acme_service = ACMEService(settings)

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "domain-svc"}

@app.post("/v1/domains/check", response_model=List[DomainAvailability])
async def check_domain_availability(domains: List[str]):
    """Check if domains are available for registration"""
    try:
        return await namecheap_service.check_availability(domains)
    except Exception as e:
        logger.error(f"Error checking domain availability: {e}")
        raise HTTPException(status_code=500, detail="Failed to check domain availability")

@app.post("/v1/domains/register", response_model=Domain)
async def register_domain(registration: DomainRegistration):
    """Register a new domain"""
    try:
        # Check if domain is available first
        availability = await namecheap_service.check_availability([registration.domain])
        if not availability[0].available:
            raise HTTPException(status_code=400, detail="Domain is not available")
        
        # Register the domain
        domain = await namecheap_service.register_domain(registration)
        
        # Set up DNS records if needed
        if registration.project_id:
            # Create A record pointing to our hosting infrastructure
            dns_record = DNSRecord(
                name=registration.domain,
                type="A",
                value="127.0.0.1",  # Replace with actual IP
                ttl=300
            )
            await dns_service.create_dns_record(dns_record)
        
        return domain
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error registering domain {registration.domain}: {e}")
        raise HTTPException(status_code=500, detail="Failed to register domain")

@app.post("/v1/domains/{domain}/dns", response_model=DNSRecord)
async def create_dns_record(domain: str, record: DNSRecord):
    """Create a DNS record for a domain"""
    try:
        return await dns_service.create_dns_record(record)
    except Exception as e:
        logger.error(f"Error creating DNS record for {domain}: {e}")
        raise HTTPException(status_code=500, detail="Failed to create DNS record")

@app.get("/v1/domains/{domain}/dns", response_model=List[DNSRecord])
async def list_dns_records(domain: str):
    """List all DNS records for a domain"""
    try:
        return await dns_service.list_dns_records(domain)
    except Exception as e:
        logger.error(f"Error listing DNS records for {domain}: {e}")
        raise HTTPException(status_code=500, detail="Failed to list DNS records")

@app.put("/v1/domains/{domain}/dns/{record_id}", response_model=DNSRecord)
async def update_dns_record(domain: str, record_id: str, record: DNSRecord):
    """Update a DNS record"""
    try:
        return await dns_service.update_dns_record(record_id, record)
    except Exception as e:
        logger.error(f"Error updating DNS record {record_id}: {e}")
        raise HTTPException(status_code=500, detail="Failed to update DNS record")

@app.delete("/v1/domains/{domain}/dns/{record_id}")
async def delete_dns_record(domain: str, record_id: str):
    """Delete a DNS record"""
    try:
        await dns_service.delete_dns_record(record_id)
        return {"message": "DNS record deleted successfully"}
    except Exception as e:
        logger.error(f"Error deleting DNS record {record_id}: {e}")
        raise HTTPException(status_code=500, detail="Failed to delete DNS record")

@app.post("/v1/domains/{domain}/certificate")
async def request_certificate(domain: str):
    """Request SSL certificate for a domain"""
    try:
        cert_info = await acme_service.request_certificate(domain, dns_service)
        if not cert_info:
            raise HTTPException(status_code=500, detail="Failed to issue certificate")
        
        return {
            "message": "Certificate issued successfully",
            "domain": domain,
            "expires_at": cert_info.expires_at
        }
    except Exception as e:
        logger.error(f"Error requesting certificate for {domain}: {e}")
        raise HTTPException(status_code=500, detail="Failed to request certificate")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=settings.port)
