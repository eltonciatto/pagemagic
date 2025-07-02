import httpx
import logging
import asyncio
from typing import List, Optional, Dict, Any
from datetime import datetime
from ..models.domain import (
    DNSRecord, DNSRecordType, Subdomain, DNSZone, 
    DomainSettings, CertificateInfo
)
from ..config.settings import Settings

logger = logging.getLogger(__name__)

class CloudflareService:
    def __init__(self, settings: Settings):
        self.api_token = settings.cloudflare_api_token
        self.zone_id = settings.cloudflare_zone_id
        self.base_url = "https://api.cloudflare.com/client/v4"
        
        self.headers = {
            "Authorization": f"Bearer {self.api_token}",
            "Content-Type": "application/json"
        }
    
    async def get_zone_id(self, domain: str) -> Optional[str]:
        """Get zone ID for a domain"""
        try:
            async with httpx.AsyncClient() as client:
                response = await client.get(
                    f"{self.base_url}/zones",
                    headers=self.headers,
                    params={"name": domain}
                )
                response.raise_for_status()
                
                result = response.json()["result"]
                if result:
                    return result[0]["id"]
                return None
                
        except Exception as e:
            logger.error(f"Error getting zone ID for {domain}: {e}")
            raise
    
    async def create_dns_record(self, zone_id: str, record: DNSRecord) -> DNSRecord:
        """Create a DNS record in Cloudflare"""
        try:
            data = {
                "type": record.type.value,
                "name": record.name,
                "content": record.value,
                "ttl": record.ttl
            }
            
            if record.priority and record.type in [DNSRecordType.MX, DNSRecordType.SRV]:
                data["priority"] = record.priority
                
            if record.weight and record.type == DNSRecordType.SRV:
                data["data"] = {
                    "weight": record.weight,
                    "port": record.port,
                    "target": record.value
                }
            
            async with httpx.AsyncClient() as client:
                response = await client.post(
                    f"{self.base_url}/zones/{zone_id}/dns_records",
                    headers=self.headers,
                    json=data
                )
                response.raise_for_status()
                
                result = response.json()["result"]
                
                return DNSRecord(
                    id=result["id"],
                    domain_id=zone_id,
                    name=result["name"],
                    type=DNSRecordType(result["type"]),
                    value=result["content"],
                    ttl=result["ttl"],
                    priority=result.get("priority"),
                    created_at=datetime.now(),
                    updated_at=datetime.now()
                )
                
        except Exception as e:
            logger.error(f"Error creating DNS record: {e}")
            raise
    
    async def update_dns_record(self, zone_id: str, record_id: str, record: DNSRecord) -> DNSRecord:
        """Update an existing DNS record"""
        try:
            data = {
                "type": record.type.value,
                "name": record.name,
                "content": record.value,
                "ttl": record.ttl
            }
            
            if record.priority and record.type in [DNSRecordType.MX, DNSRecordType.SRV]:
                data["priority"] = record.priority
            
            async with httpx.AsyncClient() as client:
                response = await client.put(
                    f"{self.base_url}/zones/{zone_id}/dns_records/{record_id}",
                    headers=self.headers,
                    json=data
                )
                response.raise_for_status()
                
                result = response.json()["result"]
                
                return DNSRecord(
                    id=result["id"],
                    domain_id=zone_id,
                    name=result["name"],
                    type=DNSRecordType(result["type"]),
                    value=result["content"],
                    ttl=result["ttl"],
                    priority=result.get("priority"),
                    updated_at=datetime.now()
                )
                
        except Exception as e:
            logger.error(f"Error updating DNS record {record_id}: {e}")
            raise
    
    async def delete_dns_record(self, zone_id: str, record_id: str) -> bool:
        """Delete a DNS record"""
        try:
            async with httpx.AsyncClient() as client:
                response = await client.delete(
                    f"{self.base_url}/zones/{zone_id}/dns_records/{record_id}",
                    headers=self.headers
                )
                response.raise_for_status()
                
                return True
                
        except Exception as e:
            logger.error(f"Error deleting DNS record {record_id}: {e}")
            raise
    
    async def list_dns_records(self, zone_id: str, name: Optional[str] = None) -> List[DNSRecord]:
        """List all DNS records for a zone"""
        try:
            params = {}
            if name:
                params["name"] = name
                
            async with httpx.AsyncClient() as client:
                response = await client.get(
                    f"{self.base_url}/zones/{zone_id}/dns_records",
                    headers=self.headers,
                    params=params
                )
                response.raise_for_status()
                
                records = []
                for record in response.json()["result"]:
                    records.append(DNSRecord(
                        id=record["id"],
                        domain_id=zone_id,
                        name=record["name"],
                        type=DNSRecordType(record["type"]),
                        value=record["content"],
                        ttl=record["ttl"],
                        priority=record.get("priority"),
                        created_at=datetime.fromisoformat(record["created_on"].replace("Z", "+00:00")),
                        updated_at=datetime.fromisoformat(record["modified_on"].replace("Z", "+00:00"))
                    ))
                
                return records
                
        except Exception as e:
            logger.error(f"Error listing DNS records for zone {zone_id}: {e}")
            raise

    async def create_subdomain(self, zone_id: str, subdomain: str, target: str, ssl_enabled: bool = False) -> Subdomain:
        """Create a subdomain with CNAME record"""
        try:
            # Create CNAME record
            record = DNSRecord(
                domain_id=zone_id,
                name=subdomain,
                type=DNSRecordType.CNAME,
                value=target,
                ttl=300
            )
            
            dns_record = await self.create_dns_record(zone_id, record)
            
            # Enable SSL if requested
            if ssl_enabled:
                await self.enable_ssl_for_subdomain(zone_id, subdomain)
            
            return Subdomain(
                id=dns_record.id,
                domain_id=zone_id,
                subdomain=subdomain,
                target=target,
                ssl_enabled=ssl_enabled,
                created_at=datetime.now(),
                updated_at=datetime.now()
            )
            
        except Exception as e:
            logger.error(f"Error creating subdomain {subdomain}: {e}")
            raise

    async def enable_ssl_for_subdomain(self, zone_id: str, hostname: str) -> bool:
        """Enable SSL for a subdomain"""
        try:
            data = {
                "type": "universal",
                "hosts": [hostname]
            }
            
            async with httpx.AsyncClient() as client:
                response = await client.post(
                    f"{self.base_url}/zones/{zone_id}/ssl/certificate_packs",
                    headers=self.headers,
                    json=data
                )
                response.raise_for_status()
                
                return True
                
        except Exception as e:
            logger.error(f"Error enabling SSL for {hostname}: {e}")
            return False

    async def get_ssl_certificate(self, zone_id: str, cert_id: str) -> Optional[CertificateInfo]:
        """Get SSL certificate information"""
        try:
            async with httpx.AsyncClient() as client:
                response = await client.get(
                    f"{self.base_url}/zones/{zone_id}/ssl/certificate_packs/{cert_id}",
                    headers=self.headers
                )
                response.raise_for_status()
                
                result = response.json()["result"]
                
                return CertificateInfo(
                    id=result["id"],
                    domain_id=zone_id,
                    domain=result["hosts"][0] if result["hosts"] else "",
                    status=result["status"],
                    expires_at=datetime.fromisoformat(result["expires_on"].replace("Z", "+00:00")) if result.get("expires_on") else None,
                    issued_at=datetime.fromisoformat(result["issued_on"].replace("Z", "+00:00")) if result.get("issued_on") else None,
                    wildcard=result.get("wildcard", False),
                    san_domains=result.get("hosts", []),
                    created_at=datetime.now(),
                    updated_at=datetime.now()
                )
                
        except Exception as e:
            logger.error(f"Error getting SSL certificate {cert_id}: {e}")
            return None

    async def purge_cache(self, zone_id: str, urls: Optional[List[str]] = None) -> bool:
        """Purge Cloudflare cache"""
        try:
            data = {}
            if urls:
                data["files"] = urls
            else:
                data["purge_everything"] = True
            
            async with httpx.AsyncClient() as client:
                response = await client.post(
                    f"{self.base_url}/zones/{zone_id}/purge_cache",
                    headers=self.headers,
                    json=data
                )
                response.raise_for_status()
                
                return True
                
        except Exception as e:
            logger.error(f"Error purging cache for zone {zone_id}: {e}")
            return False

    async def update_zone_settings(self, zone_id: str, settings: DomainSettings) -> bool:
        """Update zone settings like SSL mode, security headers, etc."""
        try:
            settings_map = {
                "ssl": "flexible" if not settings.force_https else "strict",
                "always_use_https": "on" if settings.force_https else "off",
                "automatic_https_rewrites": "on" if settings.auto_ssl else "off",
                "browser_cache_ttl": settings.cache_ttl
            }
            
            tasks = []
            for setting, value in settings_map.items():
                tasks.append(self._update_zone_setting(zone_id, setting, value))
            
            # Execute all settings updates concurrently
            await asyncio.gather(*tasks)
            
            return True
            
        except Exception as e:
            logger.error(f"Error updating zone settings for {zone_id}: {e}")
            return False

    async def _update_zone_setting(self, zone_id: str, setting: str, value: Any) -> bool:
        """Update a single zone setting"""
        try:
            data = {"value": value}
            
            async with httpx.AsyncClient() as client:
                response = await client.patch(
                    f"{self.base_url}/zones/{zone_id}/settings/{setting}",
                    headers=self.headers,
                    json=data
                )
                response.raise_for_status()
                
                return True
                
        except Exception as e:
            logger.error(f"Error updating setting {setting} for zone {zone_id}: {e}")
            return False

    async def get_zone_analytics(self, zone_id: str, since: datetime, until: datetime) -> Dict[str, Any]:
        """Get zone analytics data"""
        try:
            params = {
                "since": since.isoformat(),
                "until": until.isoformat(),
                "continuous": "true"
            }
            
            async with httpx.AsyncClient() as client:
                response = await client.get(
                    f"{self.base_url}/zones/{zone_id}/analytics/dashboard",
                    headers=self.headers,
                    params=params
                )
                response.raise_for_status()
                
                return response.json()["result"]
                
        except Exception as e:
            logger.error(f"Error getting analytics for zone {zone_id}: {e}")
            return {}
