import asyncio
import logging
import httpx
from typing import List, Optional, Dict, Any
from datetime import datetime, timedelta
from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.x509.oid import NameOID
import base64

from ..models.domain import CertificateInfo, CertificateStatus, Domain
from ..config.settings import Settings
from .acme import ACMEService

logger = logging.getLogger(__name__)

class CertificateService:
    def __init__(self, settings: Settings):
        self.settings = settings
        self.acme_service = ACMEService(settings)
        
    async def issue_certificate(self, domain: str, san_domains: List[str] = None, wildcard: bool = False) -> CertificateInfo:
        """Issue a new SSL certificate using Let's Encrypt"""
        try:
            logger.info(f"Issuing certificate for {domain}")
            
            # Generate private key
            private_key = rsa.generate_private_key(
                public_exponent=65537,
                key_size=2048
            )
            
            # Create certificate signing request
            subject = x509.Name([
                x509.NameAttribute(NameOID.COMMON_NAME, domain)
            ])
            
            csr_builder = x509.CertificateSigningRequestBuilder()
            csr_builder = csr_builder.subject_name(subject)
            
            # Add SAN (Subject Alternative Names)
            alt_names = [x509.DNSName(domain)]
            if san_domains:
                alt_names.extend([x509.DNSName(san) for san in san_domains])
            
            if wildcard:
                alt_names.append(x509.DNSName(f"*.{domain}"))
            
            csr_builder = csr_builder.add_extension(
                x509.SubjectAlternativeName(alt_names),
                critical=False
            )
            
            csr = csr_builder.sign(private_key, hashes.SHA256())
            
            # Use ACME service to get certificate
            certificate_data = await self.acme_service.request_certificate(csr, domain, san_domains)
            
            # Parse certificate to get expiry date
            cert = x509.load_pem_x509_certificate(certificate_data["certificate"].encode())
            expires_at = cert.not_valid_after
            
            return CertificateInfo(
                id=f"cert_{int(datetime.now().timestamp())}",
                domain_id="",  # Will be set by calling service
                domain=domain,
                status=CertificateStatus.ISSUED,
                certificate=certificate_data["certificate"],
                private_key=private_key.private_bytes(
                    encoding=serialization.Encoding.PEM,
                    format=serialization.PrivateFormat.PKCS8,
                    encryption_algorithm=serialization.NoEncryption()
                ).decode(),
                chain=certificate_data.get("chain", ""),
                expires_at=expires_at,
                issued_at=datetime.now(),
                wildcard=wildcard,
                san_domains=san_domains or [],
                created_at=datetime.now(),
                updated_at=datetime.now()
            )
            
        except Exception as e:
            logger.error(f"Error issuing certificate for {domain}: {e}")
            raise

    async def renew_certificate(self, certificate: CertificateInfo) -> CertificateInfo:
        """Renew an existing SSL certificate"""
        try:
            logger.info(f"Renewing certificate for {certificate.domain}")
            
            # Check if renewal is needed (within 30 days of expiry)
            if certificate.expires_at and certificate.expires_at > datetime.now() + timedelta(days=30):
                logger.info(f"Certificate for {certificate.domain} doesn't need renewal yet")
                return certificate
            
            # Issue new certificate
            new_cert = await self.issue_certificate(
                certificate.domain,
                certificate.san_domains,
                certificate.wildcard
            )
            
            # Update certificate info
            certificate.certificate = new_cert.certificate
            certificate.private_key = new_cert.private_key
            certificate.chain = new_cert.chain
            certificate.expires_at = new_cert.expires_at
            certificate.issued_at = new_cert.issued_at
            certificate.updated_at = datetime.now()
            certificate.status = CertificateStatus.ISSUED
            
            return certificate
            
        except Exception as e:
            logger.error(f"Error renewing certificate for {certificate.domain}: {e}")
            certificate.status = CertificateStatus.FAILED
            raise

    async def validate_certificate(self, certificate: CertificateInfo) -> bool:
        """Validate a certificate by checking its chain and expiry"""
        try:
            if not certificate.certificate:
                return False
            
            # Parse certificate
            cert = x509.load_pem_x509_certificate(certificate.certificate.encode())
            
            # Check if expired
            if cert.not_valid_after <= datetime.now():
                return False
            
            # Check if domain matches
            subject = cert.subject.get_attributes_for_oid(NameOID.COMMON_NAME)[0].value
            if subject != certificate.domain:
                return False
            
            # Validate certificate chain if available
            if certificate.chain:
                # TODO: Implement full chain validation
                pass
            
            return True
            
        except Exception as e:
            logger.error(f"Error validating certificate for {certificate.domain}: {e}")
            return False

    async def get_certificate_info(self, domain: str) -> Optional[Dict[str, Any]]:
        """Get certificate information for a domain via SSL check"""
        try:
            import ssl
            import socket
            
            context = ssl.create_default_context()
            
            with socket.create_connection((domain, 443), timeout=10) as sock:
                with context.wrap_socket(sock, server_hostname=domain) as ssock:
                    cert_der = ssock.getpeercert(binary_form=True)
                    cert = x509.load_der_x509_certificate(cert_der)
                    
                    return {
                        "issuer": cert.issuer.rfc4514_string(),
                        "subject": cert.subject.rfc4514_string(),
                        "expires_at": cert.not_valid_after,
                        "issued_at": cert.not_valid_before,
                        "serial_number": str(cert.serial_number),
                        "signature_algorithm": cert.signature_algorithm_oid._name
                    }
                    
        except Exception as e:
            logger.error(f"Error getting certificate info for {domain}: {e}")
            return None

    async def check_certificate_health(self, certificates: List[CertificateInfo]) -> Dict[str, Any]:
        """Check health of multiple certificates and return status"""
        results = {
            "total": len(certificates),
            "valid": 0,
            "expiring_soon": 0,
            "expired": 0,
            "failed": 0,
            "certificates": []
        }
        
        for cert in certificates:
            status = {
                "domain": cert.domain,
                "status": cert.status.value,
                "expires_at": cert.expires_at,
                "days_until_expiry": None,
                "needs_renewal": False
            }
            
            if cert.expires_at:
                days_until_expiry = (cert.expires_at - datetime.now()).days
                status["days_until_expiry"] = days_until_expiry
                
                if days_until_expiry <= 0:
                    results["expired"] += 1
                    status["needs_renewal"] = True
                elif days_until_expiry <= 30:
                    results["expiring_soon"] += 1
                    status["needs_renewal"] = True
                else:
                    results["valid"] += 1
            
            if cert.status == CertificateStatus.FAILED:
                results["failed"] += 1
            
            results["certificates"].append(status)
        
        return results

    async def auto_renew_certificates(self, certificates: List[CertificateInfo]) -> List[CertificateInfo]:
        """Automatically renew certificates that are expiring soon"""
        renewed_certificates = []
        
        for cert in certificates:
            if not cert.auto_renew:
                continue
                
            # Check if renewal is needed
            if cert.expires_at and cert.expires_at <= datetime.now() + timedelta(days=30):
                try:
                    renewed_cert = await self.renew_certificate(cert)
                    renewed_certificates.append(renewed_cert)
                    logger.info(f"Successfully renewed certificate for {cert.domain}")
                except Exception as e:
                    logger.error(f"Failed to renew certificate for {cert.domain}: {e}")
        
        return renewed_certificates

    async def revoke_certificate(self, certificate: CertificateInfo, reason: str = None) -> bool:
        """Revoke a certificate"""
        try:
            # Use ACME service to revoke certificate
            success = await self.acme_service.revoke_certificate(certificate.certificate, reason)
            
            if success:
                certificate.status = CertificateStatus.REVOKED
                certificate.updated_at = datetime.now()
                
            return success
            
        except Exception as e:
            logger.error(f"Error revoking certificate for {certificate.domain}: {e}")
            return False

    async def export_certificate(self, certificate: CertificateInfo, format: str = "pem") -> Dict[str, str]:
        """Export certificate in various formats"""
        try:
            if format.lower() == "pem":
                return {
                    "certificate": certificate.certificate,
                    "private_key": certificate.private_key,
                    "chain": certificate.chain
                }
            elif format.lower() == "p12" or format.lower() == "pfx":
                # Convert to PKCS#12 format
                from cryptography.hazmat.primitives import serialization
                
                cert = x509.load_pem_x509_certificate(certificate.certificate.encode())
                private_key = serialization.load_pem_private_key(
                    certificate.private_key.encode(),
                    password=None
                )
                
                p12 = serialization.pkcs12.serialize_key_and_certificates(
                    name=certificate.domain.encode(),
                    key=private_key,
                    cert=cert,
                    cas=None,  # Add chain certificates if needed
                    encryption_algorithm=serialization.NoEncryption()
                )
                
                return {
                    "p12": base64.b64encode(p12).decode()
                }
            else:
                raise ValueError(f"Unsupported format: {format}")
                
        except Exception as e:
            logger.error(f"Error exporting certificate for {certificate.domain}: {e}")
            raise
