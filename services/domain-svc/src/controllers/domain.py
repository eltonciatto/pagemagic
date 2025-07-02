from fastapi import APIRouter, HTTPException, Depends, status, Query
from typing import List, Optional
from datetime import datetime, timedelta

from ..models.domain import (
    Domain, DomainAvailability, DNSRecord, DomainRegistration,
    CertificateInfo, Subdomain, DomainListResponse, DNSRecordListResponse,
    CertificateListResponse, DomainSearchResult, DomainSettings
)
from ..services.dns import CloudflareService
from ..services.certificate import CertificateService
from ..services.namecheap import NamecheapService
from ..config.settings import Settings

router = APIRouter(prefix="/v1/domains", tags=["domains"])

def get_settings():
    return Settings()

def get_cloudflare_service(settings: Settings = Depends(get_settings)):
    return CloudflareService(settings)

def get_certificate_service(settings: Settings = Depends(get_settings)):
    return CertificateService(settings)

def get_namecheap_service(settings: Settings = Depends(get_settings)):
    return NamecheapService(settings)

@router.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "domain-svc"}

@router.get("/search")
async def search_domains(
    query: str = Query(..., description="Domain search query"),
    tlds: Optional[List[str]] = Query(None, description="Top-level domains to search"),
    namecheap_service: NamecheapService = Depends(get_namecheap_service)
) -> DomainSearchResult:
    """Search for available domains"""
    try:
        suggestions = await namecheap_service.search_domains(query, tlds or [".com", ".net", ".org"])
        alternatives = await namecheap_service.get_domain_suggestions(query)
        
        return DomainSearchResult(
            suggestions=suggestions,
            alternatives=alternatives
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to search domains: {str(e)}"
        )

@router.post("/register")
async def register_domain(
    registration: DomainRegistration,
    namecheap_service: NamecheapService = Depends(get_namecheap_service)
) -> Domain:
    """Register a new domain"""
    try:
        domain = await namecheap_service.register_domain(registration)
        return domain
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to register domain: {str(e)}"
        )

@router.get("/")
async def list_domains(
    user_id: str = Query(..., description="User ID"),
    page: int = Query(1, ge=1, description="Page number"),
    per_page: int = Query(10, ge=1, le=100, description="Items per page"),
    # domain_repository: DomainRepository = Depends(get_domain_repository)
) -> DomainListResponse:
    """List user domains"""
    try:
        # TODO: Implement database repository
        # domains, total = await domain_repository.get_user_domains(user_id, page, per_page)
        
        # Mock response for now
        domains = []
        total = 0
        
        return DomainListResponse(
            domains=domains,
            total=total,
            page=page,
            per_page=per_page
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to list domains: {str(e)}"
        )

@router.get("/{domain_id}")
async def get_domain(
    domain_id: str,
    # domain_repository: DomainRepository = Depends(get_domain_repository)
) -> Domain:
    """Get domain details"""
    try:
        # TODO: Implement database repository
        # domain = await domain_repository.get_domain(domain_id)
        # if not domain:
        #     raise HTTPException(status_code=404, detail="Domain not found")
        # return domain
        
        raise HTTPException(status_code=404, detail="Domain not found")
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to get domain: {str(e)}"
        )

@router.get("/{domain_id}/dns")
async def list_dns_records(
    domain_id: str,
    cloudflare_service: CloudflareService = Depends(get_cloudflare_service)
) -> DNSRecordListResponse:
    """List DNS records for a domain"""
    try:
        # Get zone ID for domain
        zone_id = await cloudflare_service.get_zone_id(domain_id)
        if not zone_id:
            raise HTTPException(status_code=404, detail="DNS zone not found")
        
        records = await cloudflare_service.list_dns_records(zone_id)
        
        return DNSRecordListResponse(
            records=records,
            domain=domain_id,
            total=len(records)
        )
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to list DNS records: {str(e)}"
        )

@router.post("/{domain_id}/dns")
async def create_dns_record(
    domain_id: str,
    record: DNSRecord,
    cloudflare_service: CloudflareService = Depends(get_cloudflare_service)
) -> DNSRecord:
    """Create a DNS record"""
    try:
        zone_id = await cloudflare_service.get_zone_id(domain_id)
        if not zone_id:
            raise HTTPException(status_code=404, detail="DNS zone not found")
        
        created_record = await cloudflare_service.create_dns_record(zone_id, record)
        return created_record
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to create DNS record: {str(e)}"
        )

@router.put("/{domain_id}/dns/{record_id}")
async def update_dns_record(
    domain_id: str,
    record_id: str,
    record: DNSRecord,
    cloudflare_service: CloudflareService = Depends(get_cloudflare_service)
) -> DNSRecord:
    """Update a DNS record"""
    try:
        zone_id = await cloudflare_service.get_zone_id(domain_id)
        if not zone_id:
            raise HTTPException(status_code=404, detail="DNS zone not found")
        
        updated_record = await cloudflare_service.update_dns_record(zone_id, record_id, record)
        return updated_record
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to update DNS record: {str(e)}"
        )

@router.delete("/{domain_id}/dns/{record_id}")
async def delete_dns_record(
    domain_id: str,
    record_id: str,
    cloudflare_service: CloudflareService = Depends(get_cloudflare_service)
):
    """Delete a DNS record"""
    try:
        zone_id = await cloudflare_service.get_zone_id(domain_id)
        if not zone_id:
            raise HTTPException(status_code=404, detail="DNS zone not found")
        
        success = await cloudflare_service.delete_dns_record(zone_id, record_id)
        if not success:
            raise HTTPException(status_code=500, detail="Failed to delete DNS record")
        
        return {"message": "DNS record deleted successfully"}
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to delete DNS record: {str(e)}"
        )

@router.get("/{domain_id}/certificates")
async def list_certificates(
    domain_id: str,
    certificate_service: CertificateService = Depends(get_certificate_service)
) -> CertificateListResponse:
    """List SSL certificates for a domain"""
    try:
        # TODO: Implement database repository to get certificates
        certificates = []
        
        return CertificateListResponse(
            certificates=certificates,
            domain=domain_id,
            total=len(certificates)
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to list certificates: {str(e)}"
        )

@router.post("/{domain_id}/certificates")
async def issue_certificate(
    domain_id: str,
    san_domains: Optional[List[str]] = None,
    wildcard: bool = False,
    certificate_service: CertificateService = Depends(get_certificate_service)
) -> CertificateInfo:
    """Issue a new SSL certificate"""
    try:
        certificate = await certificate_service.issue_certificate(
            domain_id, san_domains, wildcard
        )
        
        # TODO: Save certificate to database
        
        return certificate
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to issue certificate: {str(e)}"
        )

@router.post("/{domain_id}/certificates/{cert_id}/renew")
async def renew_certificate(
    domain_id: str,
    cert_id: str,
    certificate_service: CertificateService = Depends(get_certificate_service)
) -> CertificateInfo:
    """Renew an SSL certificate"""
    try:
        # TODO: Get certificate from database
        # certificate = await certificate_repository.get_certificate(cert_id)
        # if not certificate:
        #     raise HTTPException(status_code=404, detail="Certificate not found")
        
        # renewed_certificate = await certificate_service.renew_certificate(certificate)
        
        # TODO: Update certificate in database
        
        # return renewed_certificate
        
        raise HTTPException(status_code=404, detail="Certificate not found")
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to renew certificate: {str(e)}"
        )

@router.post("/{domain_id}/subdomains")
async def create_subdomain(
    domain_id: str,
    subdomain: str,
    target: str,
    ssl_enabled: bool = False,
    cloudflare_service: CloudflareService = Depends(get_cloudflare_service)
) -> Subdomain:
    """Create a subdomain"""
    try:
        zone_id = await cloudflare_service.get_zone_id(domain_id)
        if not zone_id:
            raise HTTPException(status_code=404, detail="DNS zone not found")
        
        subdomain_obj = await cloudflare_service.create_subdomain(
            zone_id, subdomain, target, ssl_enabled
        )
        
        # TODO: Save subdomain to database
        
        return subdomain_obj
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to create subdomain: {str(e)}"
        )

@router.post("/{domain_id}/cache/purge")
async def purge_cache(
    domain_id: str,
    urls: Optional[List[str]] = None,
    cloudflare_service: CloudflareService = Depends(get_cloudflare_service)
):
    """Purge domain cache"""
    try:
        zone_id = await cloudflare_service.get_zone_id(domain_id)
        if not zone_id:
            raise HTTPException(status_code=404, detail="DNS zone not found")
        
        success = await cloudflare_service.purge_cache(zone_id, urls)
        if not success:
            raise HTTPException(status_code=500, detail="Failed to purge cache")
        
        return {"message": "Cache purged successfully"}
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to purge cache: {str(e)}"
        )

@router.put("/{domain_id}/settings")
async def update_domain_settings(
    domain_id: str,
    settings: DomainSettings,
    cloudflare_service: CloudflareService = Depends(get_cloudflare_service)
) -> DomainSettings:
    """Update domain settings"""
    try:
        zone_id = await cloudflare_service.get_zone_id(domain_id)
        if not zone_id:
            raise HTTPException(status_code=404, detail="DNS zone not found")
        
        success = await cloudflare_service.update_zone_settings(zone_id, settings)
        if not success:
            raise HTTPException(status_code=500, detail="Failed to update settings")
        
        # TODO: Save settings to database
        
        return settings
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to update domain settings: {str(e)}"
        )

@router.get("/{domain_id}/analytics")
async def get_domain_analytics(
    domain_id: str,
    since: datetime = Query(..., description="Start date"),
    until: datetime = Query(..., description="End date"),
    cloudflare_service: CloudflareService = Depends(get_cloudflare_service)
):
    """Get domain analytics"""
    try:
        zone_id = await cloudflare_service.get_zone_id(domain_id)
        if not zone_id:
            raise HTTPException(status_code=404, detail="DNS zone not found")
        
        analytics = await cloudflare_service.get_zone_analytics(zone_id, since, until)
        return analytics
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to get analytics: {str(e)}"
        )
