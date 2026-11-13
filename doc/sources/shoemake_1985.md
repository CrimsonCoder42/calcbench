# Shoemake

**Citation**: Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242.

**Taken**: §2.3 p. 246: the quaternion product [s₁,v₁][s₂,v₂] = [s₁s₂ − v₁·v₂, s₁v₂ + s₂v₁ + v₁×v₂] (Hamilton) and addition; §3.2 pp. 247–248: rotation by conjugation x′ = q⁻¹xq (frame form) and the inverse; §3.3 p. 248: Slerp.

**Used by**: NUM-LINALG-101, NUM-LINALG-102, NUM-LINALG-103, NUM-LINALG-107, NUM-LINALG-111, NUM-LINALG-112, NUM-LINALG-113.

**Checked**: Scanned PDF OCR-read 2026-11-12 (copy at personal.fd.cvut.cz); DOI confirmed by search. Tests: ij = k, ji = −k; slerp endpoints, constant angular rate.
