from abc import ABC, abstractmethod
from typing import List, Optional, Tuple
from datetime import datetime

from ..models.domain import (
    Domain, DNSRecord, CertificateInfo, Subdomain,
    DomainRegistration, DomainSettings
)

class DomainRepositoryInterface(ABC):
    """Interface for domain repository"""
    
    @abstractmethod
    async def create_domain(self, domain: Domain) -> Domain:
        pass
    
    @abstractmethod
    async def get_domain(self, domain_id: str) -> Optional[Domain]:
        pass
    
    @abstractmethod
    async def get_domain_by_name(self, name: str) -> Optional[Domain]:
        pass
    
    @abstractmethod
    async def get_user_domains(self, user_id: str, page: int = 1, per_page: int = 10) -> Tuple[List[Domain], int]:
        pass
    
    @abstractmethod
    async def update_domain(self, domain_id: str, domain: Domain) -> Optional[Domain]:
        pass
    
    @abstractmethod
    async def delete_domain(self, domain_id: str) -> bool:
        pass

class DNSRepositoryInterface(ABC):
    """Interface for DNS records repository"""
    
    @abstractmethod
    async def create_dns_record(self, record: DNSRecord) -> DNSRecord:
        pass
    
    @abstractmethod
    async def get_dns_record(self, record_id: str) -> Optional[DNSRecord]:
        pass
    
    @abstractmethod
    async def get_domain_dns_records(self, domain_id: str) -> List[DNSRecord]:
        pass
    
    @abstractmethod
    async def update_dns_record(self, record_id: str, record: DNSRecord) -> Optional[DNSRecord]:
        pass
    
    @abstractmethod
    async def delete_dns_record(self, record_id: str) -> bool:
        pass

class CertificateRepositoryInterface(ABC):
    """Interface for certificate repository"""
    
    @abstractmethod
    async def create_certificate(self, certificate: CertificateInfo) -> CertificateInfo:
        pass
    
    @abstractmethod
    async def get_certificate(self, certificate_id: str) -> Optional[CertificateInfo]:
        pass
    
    @abstractmethod
    async def get_domain_certificates(self, domain_id: str) -> List[CertificateInfo]:
        pass
    
    @abstractmethod
    async def get_expiring_certificates(self, days: int = 30) -> List[CertificateInfo]:
        pass
    
    @abstractmethod
    async def update_certificate(self, certificate_id: str, certificate: CertificateInfo) -> Optional[CertificateInfo]:
        pass
    
    @abstractmethod
    async def delete_certificate(self, certificate_id: str) -> bool:
        pass

class SubdomainRepositoryInterface(ABC):
    """Interface for subdomain repository"""
    
    @abstractmethod
    async def create_subdomain(self, subdomain: Subdomain) -> Subdomain:
        pass
    
    @abstractmethod
    async def get_subdomain(self, subdomain_id: str) -> Optional[Subdomain]:
        pass
    
    @abstractmethod
    async def get_domain_subdomains(self, domain_id: str) -> List[Subdomain]:
        pass
    
    @abstractmethod
    async def update_subdomain(self, subdomain_id: str, subdomain: Subdomain) -> Optional[Subdomain]:
        pass
    
    @abstractmethod
    async def delete_subdomain(self, subdomain_id: str) -> bool:
        pass
