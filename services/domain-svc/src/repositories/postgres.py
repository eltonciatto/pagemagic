import asyncpg
import logging
from typing import List, Optional, Tuple, Dict, Any
from datetime import datetime, timedelta
import json

from .interfaces import (
    DomainRepositoryInterface, DNSRepositoryInterface,
    CertificateRepositoryInterface, SubdomainRepositoryInterface
)
from ..models.domain import (
    Domain, DNSRecord, CertificateInfo, Subdomain,
    DomainStatus, DNSRecordType, CertificateStatus
)

logger = logging.getLogger(__name__)

class PostgresDomainRepository(DomainRepositoryInterface):
    def __init__(self, connection_string: str):
        self.connection_string = connection_string
    
    async def _get_connection(self):
        return await asyncpg.connect(self.connection_string)
    
    async def create_domain(self, domain: Domain) -> Domain:
        conn = await self._get_connection()
        try:
            query = """
                INSERT INTO domains (id, name, status, user_id, project_id, registrar, 
                                   expires_at, auto_renew, nameservers, dns_managed, 
                                   ssl_enabled, whois_privacy, locked, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
                RETURNING *
            """
            
            row = await conn.fetchrow(
                query, domain.id, domain.name, domain.status.value, domain.user_id,
                domain.project_id, domain.registrar, domain.expires_at, domain.auto_renew,
                json.dumps(domain.nameservers), domain.dns_managed, domain.ssl_enabled,
                domain.whois_privacy, domain.locked, domain.created_at, domain.updated_at
            )
            
            return self._row_to_domain(row)
        finally:
            await conn.close()
    
    async def get_domain(self, domain_id: str) -> Optional[Domain]:
        conn = await self._get_connection()
        try:
            query = "SELECT * FROM domains WHERE id = $1"
            row = await conn.fetchrow(query, domain_id)
            
            return self._row_to_domain(row) if row else None
        finally:
            await conn.close()
    
    async def get_domain_by_name(self, name: str) -> Optional[Domain]:
        conn = await self._get_connection()
        try:
            query = "SELECT * FROM domains WHERE name = $1"
            row = await conn.fetchrow(query, name)
            
            return self._row_to_domain(row) if row else None
        finally:
            await conn.close()
    
    async def get_user_domains(self, user_id: str, page: int = 1, per_page: int = 10) -> Tuple[List[Domain], int]:
        conn = await self._get_connection()
        try:
            # Get total count
            count_query = "SELECT COUNT(*) FROM domains WHERE user_id = $1"
            total = await conn.fetchval(count_query, user_id)
            
            # Get paginated results
            offset = (page - 1) * per_page
            query = """
                SELECT * FROM domains 
                WHERE user_id = $1 
                ORDER BY created_at DESC 
                LIMIT $2 OFFSET $3
            """
            
            rows = await conn.fetch(query, user_id, per_page, offset)
            domains = [self._row_to_domain(row) for row in rows]
            
            return domains, total
        finally:
            await conn.close()
    
    async def update_domain(self, domain_id: str, domain: Domain) -> Optional[Domain]:
        conn = await self._get_connection()
        try:
            query = """
                UPDATE domains 
                SET status = $2, registrar = $3, expires_at = $4, auto_renew = $5,
                    nameservers = $6, dns_managed = $7, ssl_enabled = $8, 
                    whois_privacy = $9, locked = $10, updated_at = $11
                WHERE id = $1
                RETURNING *
            """
            
            row = await conn.fetchrow(
                query, domain_id, domain.status.value, domain.registrar,
                domain.expires_at, domain.auto_renew, json.dumps(domain.nameservers),
                domain.dns_managed, domain.ssl_enabled, domain.whois_privacy,
                domain.locked, datetime.now()
            )
            
            return self._row_to_domain(row) if row else None
        finally:
            await conn.close()
    
    async def delete_domain(self, domain_id: str) -> bool:
        conn = await self._get_connection()
        try:
            query = "DELETE FROM domains WHERE id = $1"
            result = await conn.execute(query, domain_id)
            
            return result == "DELETE 1"
        finally:
            await conn.close()
    
    def _row_to_domain(self, row) -> Domain:
        nameservers = json.loads(row["nameservers"]) if row["nameservers"] else []
        
        return Domain(
            id=row["id"],
            name=row["name"],
            status=DomainStatus(row["status"]),
            user_id=row["user_id"],
            project_id=row["project_id"],
            registrar=row["registrar"],
            expires_at=row["expires_at"],
            auto_renew=row["auto_renew"],
            nameservers=nameservers,
            dns_managed=row["dns_managed"],
            ssl_enabled=row["ssl_enabled"],
            whois_privacy=row["whois_privacy"],
            locked=row["locked"],
            created_at=row["created_at"],
            updated_at=row["updated_at"]
        )

class PostgresDNSRepository(DNSRepositoryInterface):
    def __init__(self, connection_string: str):
        self.connection_string = connection_string
    
    async def _get_connection(self):
        return await asyncpg.connect(self.connection_string)
    
    async def create_dns_record(self, record: DNSRecord) -> DNSRecord:
        conn = await self._get_connection()
        try:
            query = """
                INSERT INTO dns_records (id, domain_id, name, type, value, ttl, priority, 
                                       weight, port, enabled, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
                RETURNING *
            """
            
            row = await conn.fetchrow(
                query, record.id, record.domain_id, record.name, record.type.value,
                record.value, record.ttl, record.priority, record.weight, record.port,
                record.enabled, record.created_at or datetime.now(), 
                record.updated_at or datetime.now()
            )
            
            return self._row_to_dns_record(row)
        finally:
            await conn.close()
    
    async def get_dns_record(self, record_id: str) -> Optional[DNSRecord]:
        conn = await self._get_connection()
        try:
            query = "SELECT * FROM dns_records WHERE id = $1"
            row = await conn.fetchrow(query, record_id)
            
            return self._row_to_dns_record(row) if row else None
        finally:
            await conn.close()
    
    async def get_domain_dns_records(self, domain_id: str) -> List[DNSRecord]:
        conn = await self._get_connection()
        try:
            query = "SELECT * FROM dns_records WHERE domain_id = $1 ORDER BY name, type"
            rows = await conn.fetch(query, domain_id)
            
            return [self._row_to_dns_record(row) for row in rows]
        finally:
            await conn.close()
    
    async def update_dns_record(self, record_id: str, record: DNSRecord) -> Optional[DNSRecord]:
        conn = await self._get_connection()
        try:
            query = """
                UPDATE dns_records 
                SET name = $2, type = $3, value = $4, ttl = $5, priority = $6,
                    weight = $7, port = $8, enabled = $9, updated_at = $10
                WHERE id = $1
                RETURNING *
            """
            
            row = await conn.fetchrow(
                query, record_id, record.name, record.type.value, record.value,
                record.ttl, record.priority, record.weight, record.port,
                record.enabled, datetime.now()
            )
            
            return self._row_to_dns_record(row) if row else None
        finally:
            await conn.close()
    
    async def delete_dns_record(self, record_id: str) -> bool:
        conn = await self._get_connection()
        try:
            query = "DELETE FROM dns_records WHERE id = $1"
            result = await conn.execute(query, record_id)
            
            return result == "DELETE 1"
        finally:
            await conn.close()
    
    def _row_to_dns_record(self, row) -> DNSRecord:
        return DNSRecord(
            id=row["id"],
            domain_id=row["domain_id"],
            name=row["name"],
            type=DNSRecordType(row["type"]),
            value=row["value"],
            ttl=row["ttl"],
            priority=row["priority"],
            weight=row["weight"],
            port=row["port"],
            enabled=row["enabled"],
            created_at=row["created_at"],
            updated_at=row["updated_at"]
        )

class PostgresCertificateRepository(CertificateRepositoryInterface):
    def __init__(self, connection_string: str):
        self.connection_string = connection_string
    
    async def _get_connection(self):
        return await asyncpg.connect(self.connection_string)
    
    async def create_certificate(self, certificate: CertificateInfo) -> CertificateInfo:
        conn = await self._get_connection()
        try:
            query = """
                INSERT INTO certificates (id, domain_id, domain, status, certificate, 
                                        private_key, chain, expires_at, issued_at, 
                                        auto_renew, wildcard, san_domains, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
                RETURNING *
            """
            
            row = await conn.fetchrow(
                query, certificate.id, certificate.domain_id, certificate.domain,
                certificate.status.value, certificate.certificate, certificate.private_key,
                certificate.chain, certificate.expires_at, certificate.issued_at,
                certificate.auto_renew, certificate.wildcard, 
                json.dumps(certificate.san_domains), certificate.created_at, certificate.updated_at
            )
            
            return self._row_to_certificate(row)
        finally:
            await conn.close()
    
    async def get_certificate(self, certificate_id: str) -> Optional[CertificateInfo]:
        conn = await self._get_connection()
        try:
            query = "SELECT * FROM certificates WHERE id = $1"
            row = await conn.fetchrow(query, certificate_id)
            
            return self._row_to_certificate(row) if row else None
        finally:
            await conn.close()
    
    async def get_domain_certificates(self, domain_id: str) -> List[CertificateInfo]:
        conn = await self._get_connection()
        try:
            query = "SELECT * FROM certificates WHERE domain_id = $1 ORDER BY created_at DESC"
            rows = await conn.fetch(query, domain_id)
            
            return [self._row_to_certificate(row) for row in rows]
        finally:
            await conn.close()
    
    async def get_expiring_certificates(self, days: int = 30) -> List[CertificateInfo]:
        conn = await self._get_connection()
        try:
            expiry_threshold = datetime.now() + timedelta(days=days)
            query = """
                SELECT * FROM certificates 
                WHERE expires_at <= $1 AND status = 'issued' AND auto_renew = true
                ORDER BY expires_at ASC
            """
            rows = await conn.fetch(query, expiry_threshold)
            
            return [self._row_to_certificate(row) for row in rows]
        finally:
            await conn.close()
    
    async def update_certificate(self, certificate_id: str, certificate: CertificateInfo) -> Optional[CertificateInfo]:
        conn = await self._get_connection()
        try:
            query = """
                UPDATE certificates 
                SET status = $2, certificate = $3, private_key = $4, chain = $5,
                    expires_at = $6, issued_at = $7, auto_renew = $8, updated_at = $9
                WHERE id = $1
                RETURNING *
            """
            
            row = await conn.fetchrow(
                query, certificate_id, certificate.status.value, certificate.certificate,
                certificate.private_key, certificate.chain, certificate.expires_at,
                certificate.issued_at, certificate.auto_renew, datetime.now()
            )
            
            return self._row_to_certificate(row) if row else None
        finally:
            await conn.close()
    
    async def delete_certificate(self, certificate_id: str) -> bool:
        conn = await self._get_connection()
        try:
            query = "DELETE FROM certificates WHERE id = $1"
            result = await conn.execute(query, certificate_id)
            
            return result == "DELETE 1"
        finally:
            await conn.close()
    
    def _row_to_certificate(self, row) -> CertificateInfo:
        san_domains = json.loads(row["san_domains"]) if row["san_domains"] else []
        
        return CertificateInfo(
            id=row["id"],
            domain_id=row["domain_id"],
            domain=row["domain"],
            status=CertificateStatus(row["status"]),
            certificate=row["certificate"],
            private_key=row["private_key"],
            chain=row["chain"],
            expires_at=row["expires_at"],
            issued_at=row["issued_at"],
            auto_renew=row["auto_renew"],
            wildcard=row["wildcard"],
            san_domains=san_domains,
            created_at=row["created_at"],
            updated_at=row["updated_at"]
        )
