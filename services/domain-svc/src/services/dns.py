import httpx
import logging
from typing import List, Optional
from ..models.domain import DNSRecord, DNSRecordType
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
    
    async def create_dns_record(self, record: DNSRecord) -> DNSRecord:
        """Create a DNS record in Cloudflare"""
        try:
            data = {
                "type": record.type,
                "name": record.name,
                "content": record.value,
                "ttl": record.ttl
            }
            
            if record.priority:
                data["priority"] = record.priority
            
            async with httpx.AsyncClient() as client:
                response = await client.post(
                    f"{self.base_url}/zones/{self.zone_id}/dns_records",
                    headers=self.headers,
                    json=data
                )
                response.raise_for_status()
                
                result = response.json()["result"]
                
                return DNSRecord(
                    id=result["id"],
                    name=result["name"],
                    type=result["type"],
                    value=result["content"],
                    ttl=result["ttl"],
                    priority=result.get("priority")
                )
                
        except Exception as e:
            logger.error(f"Error creating DNS record: {e}")
            raise
    
    async def update_dns_record(self, record_id: str, record: DNSRecord) -> DNSRecord:
        """Update an existing DNS record"""
        try:
            data = {
                "type": record.type,
                "name": record.name,
                "content": record.value,
                "ttl": record.ttl
            }
            
            if record.priority:
                data["priority"] = record.priority
            
            async with httpx.AsyncClient() as client:
                response = await client.put(
                    f"{self.base_url}/zones/{self.zone_id}/dns_records/{record_id}",
                    headers=self.headers,
                    json=data
                )
                response.raise_for_status()
                
                result = response.json()["result"]
                
                return DNSRecord(
                    id=result["id"],
                    name=result["name"],
                    type=result["type"],
                    value=result["content"],
                    ttl=result["ttl"],
                    priority=result.get("priority")
                )
                
        except Exception as e:
            logger.error(f"Error updating DNS record {record_id}: {e}")
            raise
    
    async def delete_dns_record(self, record_id: str) -> bool:
        """Delete a DNS record"""
        try:
            async with httpx.AsyncClient() as client:
                response = await client.delete(
                    f"{self.base_url}/zones/{self.zone_id}/dns_records/{record_id}",
                    headers=self.headers
                )
                response.raise_for_status()
                
                return True
                
        except Exception as e:
            logger.error(f"Error deleting DNS record {record_id}: {e}")
            raise
    
    async def list_dns_records(self, domain: str) -> List[DNSRecord]:
        """List all DNS records for a domain"""
        try:
            async with httpx.AsyncClient() as client:
                response = await client.get(
                    f"{self.base_url}/zones/{self.zone_id}/dns_records",
                    headers=self.headers,
                    params={"name": domain}
                )
                response.raise_for_status()
                
                records = []
                for record in response.json()["result"]:
                    records.append(DNSRecord(
                        id=record["id"],
                        name=record["name"],
                        type=record["type"],
                        value=record["content"],
                        ttl=record["ttl"],
                        priority=record.get("priority")
                    ))
                
                return records
                
        except Exception as e:
            logger.error(f"Error listing DNS records for {domain}: {e}")
            raise
