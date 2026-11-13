# Shepperd

**Citation**: Shepperd, "Quaternion from Rotation Matrix", Journal of Guidance and Control, 1978, 1(3):223–224, doi:10.2514/3.55767b.

**Taken**: Choose the largest of 1 + trace and 1 + 2mᵢᵢ − trace before the square root, so no division is by a small number.

**Used by**: NUM-LINALG-110.

**Checked**: Citation confirmed by search 2026-11-12 (Crossref rate-limited); formula cross-checked with Diebel §6.5 eqs. (141)–(145). Test: FromMat3(ToMat3(q)) = ±q on 500 random rotations.
