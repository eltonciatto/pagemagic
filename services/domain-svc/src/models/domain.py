from pydantic import BaseModel
from typing import Optional, List
from datetime import datetime
from enum import Enum

class DomainStatus(str, Enum):
    AVAILABLE = "available"
    REGISTERED = "registered" 
    PENDING = "pending"
    EXPIRED = "expired"
    TRANSFERRED = "transferred"

class DNSRecordType(str, Enum):
    A = "A"
    AAAA = "AAAA"
    CNAME = "CNAME"
    MX = "MX"
    TXT = "TXT"

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
    created_at: datetime
    updated_at: datetime

class DomainAvailability(BaseModel):
    domain: str
    available: bool
    price: Optional[float] = None
    currency: str = "USD"
    premium: bool = False

class DNSRecord(BaseModel):
    id: Optional[str] = None
    name: str
    type: DNSRecordType
    value: str
    ttl: int = 300
    priority: Optional[int] = None

class DomainRegistration(BaseModel):
    domain: str
    user_id: str
    project_id: Optional[str] = None
    years: int = 1
    auto_renew: bool = True
    registrant_info: dict
    nameservers: Optional[List[str]] = None

class CertificateInfo(BaseModel):
    domain: str
    certificate: str
    private_key: str
    chain: str
    expires_at: datetime
    issued_at: datetime
