# Diebel

**Citation**: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006.

**Taken**: §2.4 eqs. (14)–(16) single-axis rotations; §5.6 the (1, 2, 3) sequence — roll φ, pitch θ, yaw ψ — eq. (67) matrix, eqs. (72)–(73) inverse, eq. (84) quaternion, §5.6.6 and §5.8 singularities (gimbal lock); §6.1 eqs. (98)–(100); §6.4 eqs. (120)–(127); §6.5 eqs. (131)–(145); §6.12 eq. (175); §6.13 eqs. (198)–(199); §7.7 eqs. (243)–(244). Diebel's matrices are FRAME (passive) rotations and his quaternion product carries −q×p; the active matrices here are his transposes, and the quaternion components of the (1, 2, 3) sequence coincide.

**Used by**: NUM-LINALG-010, NUM-LINALG-031, NUM-LINALG-063, NUM-LINALG-072, NUM-LINALG-073, NUM-LINALG-074, NUM-LINALG-075, NUM-LINALG-104, NUM-LINALG-105, NUM-LINALG-106, NUM-LINALG-107, NUM-LINALG-108, NUM-LINALG-109, NUM-LINALG-110, NUM-LINALG-111, NUM-LINALG-112, NUM-LINALG-114, NUM-LINALG-115.

**Checked**: PDF (copy at https://www.astro.rug.nl/software/kapteyn-beta/_downloads/attitude.pdf) text extracted and read 2026-11-12. Tests write eq. (14)–(16), eq. (67) and eq. (125) out literally and compare (TestMat3Rotation_diebelTransposes, TestQuat_oracle).
