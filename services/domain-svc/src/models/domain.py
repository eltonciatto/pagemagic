from pydantic import BaseModel, validator
from typing import Optional, List, Dict, Any
from datetime import datetime
from enum import Enum

class DomainStatus(str, Enum):
    AVAILABLE = "available"
    REGISTERED = "registered" 
    PENDING = "pending"
    EXPIRED = "expired"
    TRANSFERRED = "transferred"
    SUSPENDED = "suspended"
    LOCKED = "locked"

class DNSRecordType(str, Enum):
    A = "A"
    AAAA = "AAAA"
    CNAME = "CNAME"
    MX = "MX"
    TXT = "TXT"
    SRV = "SRV"
    NS = "NS"
    PTR = "PTR"

class CertificateStatus(str, Enum):
    PENDING = "pending"
    ISSUED = "issued"
    EXPIRED = "expired"
    REVOKED = "revoked"
    FAILED = "failed"

class Domain(BaseModel):
    id: str
    name: str
    status: DomainStatus
    user_id: str
    project_id: Optional[str] = None
    registrar: Optional[str] = None
    expires_at: Optional[datetime] = None
    auto_renew: bool = True
    nameservers: List[str] = []
    dns_managed: bool = False
    ssl_enabled: bool = False
    whois_privacy: bool = True
    locked: bool = False
    created_at: datetime
    updated_at: datetime

    @validator('name')
    def validate_domain_name(cls, v):
        if not v or '.' not in v:
            raise ValueError('Invalid domain name')
        return v.lower()

class Subdomain(BaseModel):
    id: str
    domain_id: str
    subdomain: str
    target: str
    enabled: bool = True
    ssl_enabled: bool = False
    created_at: datetime
    updated_at: datetime

class DomainAvailability(BaseModel):
    domain: str
    available: bool
    price: Optional[float] = None
    currency: str = "USD"
    premium: bool = False
    renewal_price: Optional[float] = None
    transfer_price: Optional[float] = None

class DNSRecord(BaseModel):
    id: Optional[str] = None
    domain_id: str
    name: str
    type: DNSRecordType
    value: str
    ttl: int = 300
    priority: Optional[int] = None
    weight: Optional[int] = None
    port: Optional[int] = None
    enabled: bool = True
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None

    @validator('ttl')
    def validate_ttl(cls, v):
        if v < 60 or v > 86400:
            raise ValueError('TTL must be between 60 and 86400 seconds')
        return v

class DomainRegistration(BaseModel):
    domain: str
    user_id: str
    project_id: Optional[str] = None
    years: int = 1
    auto_renew: bool = True
    registrant_info: Dict[str, Any]
    nameservers: Optional[List[str]] = None
    whois_privacy: bool = True

    @validator('years')
    def validate_years(cls, v):
        if v < 1 or v > 10:
            raise ValueError('Years must be between 1 and 10')
        return v

class RegistrantInfo(BaseModel):
    first_name: str
    last_name: str
    organization: Optional[str] = None
    email: str
    phone: str
    address1: str
    address2: Optional[str] = None
    city: str
    state: str
    postal_code: str
    country: str

class CertificateInfo(BaseModel):
    id: str
    domain_id: str
    domain: str
    status: CertificateStatus
    certificate: Optional[str] = None
    private_key: Optional[str] = None
    chain: Optional[str] = None
    expires_at: Optional[datetime] = None
    issued_at: Optional[datetime] = None
    auto_renew: bool = True
    wildcard: bool = False
    san_domains: List[str] = []
    created_at: datetime
    updated_at: datetime

class DomainTransfer(BaseModel):
    id: str
    domain: str
    user_id: str
    auth_code: str
    status: str
    registrant_info: Dict[str, Any]
    created_at: datetime
    updated_at: datetime

class DNSZone(BaseModel):
    id: str
    domain_id: str
    zone_file: str
    serial: int
    refresh: int = 3600
    retry: int = 1800
    expire: int = 604800
    minimum: int = 86400
    created_at: datetime
    updated_at: datetime

class DomainSettings(BaseModel):
    domain_id: str
    auto_ssl: bool = True
    force_https: bool = False
    cdn_enabled: bool = False
    cache_ttl: int = 3600
    custom_error_pages: Dict[str, str] = {}
    security_headers: Dict[str, str] = {}

# API Response Models
class DomainListResponse(BaseModel):
    domains: List[Domain]
    total: int
    page: int
    per_page: int

class DNSRecordListResponse(BaseModel):
    records: List[DNSRecord]
    domain: str
    total: int

class CertificateListResponse(BaseModel):
    certificates: List[CertificateInfo]
    domain: str
    total: int

class DomainSearchResult(BaseModel):
    suggestions: List[DomainAvailability]
    alternatives: List[DomainAvailability]

# Event Models for messaging
class DomainEvent(BaseModel):
    event_id: str
    event_type: str
    domain_id: str
    user_id: str
    data: Dict[str, Any]
    timestamp: datetime

class CertificateRenewalEvent(BaseModel):
    certificate_id: str
    domain: str
    expires_at: datetime
    auto_renew: bool
    days_until_expiry: int
