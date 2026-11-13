// Package linalg holds vectors, matrices, quaternions and 3D coordinate conversions, and the standard matrix
// toolkit (elementwise and Kronecker products, echelon forms and subspaces, cofactors, checks, Cholesky, QR,
// least squares, triangular and tridiagonal solves, symmetric and 2×2 eigenproblems, SVD, norms): fixed-size Vec2,
// Vec3, Vec4, Mat2, Mat3, Mat4 and Quat values that never allocate, and n-dimensional functions on
// []float64 (vectors, and row-major matrices with their dimensions passed alongside).
//
// ID short form: LINALG (IDs are NUM-LINALG-NNN).
//
// Conventions. float64 throughout; angles in radians. Vectors are columns, and rotations are ACTIVE
// (they turn the vector, counterclockwise about the axis by the right-hand rule); Diebel 2006, cited
// below, writes frame rotations, the transposes. Quaternions are Hamilton's, scalar part W first by name.
// Euler angles are aerospace yaw-pitch-roll (Z-Y-X). Spherical coordinates are the physics convention of
// DLMF §1.5(ii): θ from +z, φ azimuth in [0, 2π). No function modifies its inputs. Functions that produce
// an n-dimensional result write into a dst slice, reused when its capacity suffices and allocated once
// otherwise; dst must not overlap an input unless the function says it may. Tolerances are parameters
// with a documented default. NaN or ±Inf in returns ErrNotFinite; finite inputs with a result outside the
// float64 range return ErrOverflow; a zero vector where a direction is needed returns ErrZeroVector; a
// singular matrix returns ErrSingular — never a silent NaN. Errors are wrapped with the parameter name;
// test with errors.Is. 2D points, distances and polar coordinates live in num/geometry; angle
// normalisation lives in num/trig.
//
// Functions:
//
//	NUM-LINALG-001  Vec2Add                   a + b in the plane
//	NUM-LINALG-002  Vec2Sub                   a − b in the plane
//	NUM-LINALG-003  Vec2Scale                 k·v in the plane
//	NUM-LINALG-004  Vec2Div                   v/k in the plane, error at k = 0
//	NUM-LINALG-005  Vec2Dot                   dot product of two 2D vectors, within 2 ulps
//	NUM-LINALG-006  Vec2Cross                 2D scalar cross product (determinant, right-hand-rule sign)
//	NUM-LINALG-007  Vec2Norm                  length of a 2D vector without overflow
//	NUM-LINALG-008  Vec2Normalize             unit vector of a 2D vector, error at zero
//	NUM-LINALG-009  Vec2AngleBetween          angle between two 2D vectors (Kahan's formula)
//	NUM-LINALG-010  Vec2Rotate                2D vector turned counterclockwise by an angle
//	NUM-LINALG-011  Vec3Add                   a + b in space
//	NUM-LINALG-012  Vec3Sub                   a − b in space
//	NUM-LINALG-013  Vec3Scale                 k·v in space
//	NUM-LINALG-014  Vec3Div                   v/k in space, error at k = 0
//	NUM-LINALG-015  Vec3Dot                   dot product of two 3D vectors
//	NUM-LINALG-016  Vec3Cross                 cross product, each component within 2 ulps
//	NUM-LINALG-017  Vec3Norm                  length of a 3D vector without overflow
//	NUM-LINALG-018  Vec3Normalize             unit vector of a 3D vector, error at zero
//	NUM-LINALG-019  Vec3Distance              distance between two points in space
//	NUM-LINALG-020  Vec3AngleBetween          angle between two 3D vectors (Kahan's formula)
//	NUM-LINALG-021  Vec3DirectionAngles       angles between a vector and the x, y, z axes
//	NUM-LINALG-022  Vec3ComponentAlong        scalar component of v along an axis
//	NUM-LINALG-023  Vec3Project               projection of v onto a line
//	NUM-LINALG-024  Vec3Reject                component of v perpendicular to a line
//	NUM-LINALG-025  Vec3Reflect               mirror image of v in a plane through the origin
//	NUM-LINALG-026  Vec3ParallelogramArea     area of the parallelogram spanned by two vectors
//	NUM-LINALG-027  Vec3TriangleArea          area of the triangle spanned by two vectors
//	NUM-LINALG-028  Vec3ScalarTripleProduct   a·(b × c), the signed volume (mixed product)
//	NUM-LINALG-029  Vec3VectorTripleProduct   a × (b × c)
//	NUM-LINALG-030  Vec3PlaneNormal           unit normal of the plane through three points
//	NUM-LINALG-031  Vec3Rotate                v turned about an axis by an angle (Rodrigues)
//	NUM-LINALG-032  Vec3Lerp                  linear interpolation, exact at both ends
//	NUM-LINALG-033  Vec3Min                   componentwise minimum
//	NUM-LINALG-034  Vec3Max                   componentwise maximum
//	NUM-LINALG-035  VecAdd                    a + b in n dimensions
//	NUM-LINALG-036  VecSub                    a − b in n dimensions
//	NUM-LINALG-037  VecScale                  k·v in n dimensions
//	NUM-LINALG-038  VecDiv                    v/k in n dimensions, error at k = 0
//	NUM-LINALG-039  VecDot                    dot product in n dimensions
//	NUM-LINALG-040  VecDotCompensated         dot product as if in twice the working precision (Dot2)
//	NUM-LINALG-041  VecNorm                   Euclidean norm, scaled and compensated
//	NUM-LINALG-042  VecNorm1                  1-norm, compensated sum
//	NUM-LINALG-043  VecNormInf                infinity norm (largest |element|)
//	NUM-LINALG-044  VecNormP                  p-norm for p ≥ 1, scaled
//	NUM-LINALG-045  VecNormalize              unit vector in n dimensions, error at zero
//	NUM-LINALG-046  VecDistance               Euclidean distance between two points in n dimensions
//	NUM-LINALG-047  VecAngleBetween           angle between two n-D vectors (Kahan's formula)
//	NUM-LINALG-048  VecProject                projection of v onto a line in n dimensions
//	NUM-LINALG-049  VecLerp                   linear interpolation in n dimensions, exact at both ends
//	NUM-LINALG-050  VecMin                    element-by-element minimum
//	NUM-LINALG-051  VecMax                    element-by-element maximum
//	NUM-LINALG-052  CartesianToCylindrical    (x, y, z) to cylindrical (R, φ, z)
//	NUM-LINALG-053  CylindricalToCartesian    cylindrical (R, φ, z) to (x, y, z)
//	NUM-LINALG-054  CartesianToSpherical      (x, y, z) to spherical (R, θ polar, φ azimuth)
//	NUM-LINALG-055  SphericalToCartesian      spherical (R, θ, φ) to (x, y, z)
//	NUM-LINALG-056  Mat2Identity              2×2 identity
//	NUM-LINALG-057  Mat2Mul                   2×2 matrix product, elements within 2 ulps
//	NUM-LINALG-058  Mat2MulVec                2×2 matrix times a 2D vector
//	NUM-LINALG-059  Mat2Transpose             2×2 transpose
//	NUM-LINALG-060  Mat2Det                   2×2 determinant within 2 ulps
//	NUM-LINALG-061  Mat2Inverse               2×2 inverse, error when singular
//	NUM-LINALG-062  Mat2Trace                 2×2 trace
//	NUM-LINALG-063  Mat2Rotation              2×2 rotation matrix
//	NUM-LINALG-064  Mat2SymmetricEigenvalues  eigenvalues of a symmetric 2×2 matrix without cancellation
//	NUM-LINALG-065  Mat3Identity              3×3 identity
//	NUM-LINALG-066  Mat3Mul                   3×3 matrix product
//	NUM-LINALG-067  Mat3MulVec                3×3 matrix times a 3D vector
//	NUM-LINALG-068  Mat3Transpose             3×3 transpose
//	NUM-LINALG-069  Mat3Det                   3×3 determinant as a triple product
//	NUM-LINALG-070  Mat3Inverse               3×3 inverse by cofactors, error when singular
//	NUM-LINALG-071  Mat3Trace                 3×3 trace
//	NUM-LINALG-072  Mat3RotationX             rotation matrix about the x axis
//	NUM-LINALG-073  Mat3RotationY             rotation matrix about the y axis
//	NUM-LINALG-074  Mat3RotationZ             rotation matrix about the z axis
//	NUM-LINALG-075  Mat3RotationAxisAngle     rotation matrix about any axis (Rodrigues)
//	NUM-LINALG-076  Mat3SymmetricEigenvalues  eigenvalues of a symmetric 3×3 matrix (Smith's closed form)
//	NUM-LINALG-077  Mat4Identity              4×4 identity
//	NUM-LINALG-078  Mat4Mul                   4×4 matrix product
//	NUM-LINALG-079  Mat4MulVec                4×4 matrix times a 4-vector
//	NUM-LINALG-080  Mat4Transpose             4×4 transpose
//	NUM-LINALG-081  Mat4Det                   4×4 determinant by Laplace expansion
//	NUM-LINALG-082  Mat4Inverse               4×4 inverse by cofactors, error when singular
//	NUM-LINALG-083  Mat4Trace                 4×4 trace
//	NUM-LINALG-084  Mat4Translation           homogeneous translation matrix
//	NUM-LINALG-085  Mat4Scale                 homogeneous scaling matrix
//	NUM-LINALG-086  Mat4FromMat3              3×3 linear map embedded in a 4×4 homogeneous matrix
//	NUM-LINALG-087  Mat4TransformPoint        point transformed by a homogeneous matrix (with the w divide)
//	NUM-LINALG-088  Mat4TransformDirection    direction transformed by a homogeneous matrix (no translation)
//	NUM-LINALG-089  MatIdentity               n×n identity
//	NUM-LINALG-090  MatMul                    product of two general matrices
//	NUM-LINALG-091  MatVecMul                 general matrix times a vector
//	NUM-LINALG-092  MatTranspose              transpose of a general matrix
//	NUM-LINALG-093  MatTrace                  trace of an n×n matrix, compensated
//	NUM-LINALG-094  LUDecompose               PA = LU with partial pivoting
//	NUM-LINALG-095  LUSolve                   solve Ax = b from an LU factorisation
//	NUM-LINALG-096  MatSolve                  solve Ax = b by LU with partial pivoting
//	NUM-LINALG-097  MatDet                    n×n determinant from LU, without intermediate overflow
//	NUM-LINALG-098  MatInverse                n×n inverse from LU, error when singular
//	NUM-LINALG-099  MatRank                   numerical rank by complete pivoting with a tolerance
//	NUM-LINALG-100  MatConditionEstimate1     estimate of the 1-norm condition number (Hager–Higham)
//	NUM-LINALG-101  QuatAdd                   p + q for quaternions
//	NUM-LINALG-102  QuatSub                   p − q for quaternions
//	NUM-LINALG-103  QuatMul                   Hamilton product p·q
//	NUM-LINALG-104  QuatConjugate             quaternion conjugate
//	NUM-LINALG-105  QuatNorm                  quaternion magnitude without overflow
//	NUM-LINALG-106  QuatNormalize             unit quaternion (versor), error at zero
//	NUM-LINALG-107  QuatInverse               quaternion inverse, error at zero
//	NUM-LINALG-108  QuatFromAxisAngle         unit quaternion of a rotation about an axis
//	NUM-LINALG-109  QuatToAxisAngle           axis and angle of a quaternion's rotation
//	NUM-LINALG-110  QuatFromMat3              unit quaternion of a rotation matrix (Shepperd)
//	NUM-LINALG-111  QuatToMat3                rotation matrix of a quaternion
//	NUM-LINALG-112  QuatRotateVec3            3D vector rotated by a quaternion
//	NUM-LINALG-113  QuatSlerp                 spherical linear interpolation between rotations
//	NUM-LINALG-114  QuatFromEulerZYX          quaternion from yaw, pitch, roll (aerospace Z-Y-X)
//	NUM-LINALG-115  QuatToEulerZYX            yaw, pitch, roll from a quaternion, gimbal lock handled
//	NUM-LINALG-116  MatAdd                       A + B entry by entry
//	NUM-LINALG-117  MatSub                       A − B entry by entry
//	NUM-LINALG-118  MatScale                     k·A
//	NUM-LINALG-119  MatHadamard                  Hadamard (entrywise) product A ∘ B
//	NUM-LINALG-120  VecOuter                     outer product u·vᵀ
//	NUM-LINALG-121  MatKronecker                 Kronecker product A ⊗ B
//	NUM-LINALG-122  MatPower                     Aᵏ by repeated squaring
//	NUM-LINALG-123  Mat3Minor                    3×3 minor: 2×2 determinant with row i, column j deleted
//	NUM-LINALG-124  Mat2Adjugate                 2×2 adjugate [d −b; −c a]
//	NUM-LINALG-125  Mat3Adjugate                 3×3 adjugate (transposed cofactors)
//	NUM-LINALG-126  Mat4Adjugate                 4×4 adjugate (transposed cofactors)
//	NUM-LINALG-127  MatMinor                     n×n minor by LU
//	NUM-LINALG-128  MatCofactorMatrix            n×n cofactor matrix (det·A⁻ᵀ, or minors when singular)
//	NUM-LINALG-129  MatAdjugate                  n×n adjugate, A·adj A = det(A)·I
//	NUM-LINALG-130  MatRowEchelon                row echelon form with partial pivoting and tolerance
//	NUM-LINALG-131  MatRREF                      reduced row echelon form (Gauss–Jordan)
//	NUM-LINALG-132  MatNullSpace                 null-space basis from the RREF
//	NUM-LINALG-133  MatColumnSpace               column-space basis: the pivot columns of A
//	NUM-LINALG-134  Mat2Characteristics          det, trace, rank, inverse, eigenvalues of a 2×2
//	NUM-LINALG-135  Mat3Characteristics          det, trace, rank, inverse, char. polynomial, eigenvalues of a 3×3
//	NUM-LINALG-136  MatIsSymmetric               symmetric to a tolerance
//	NUM-LINALG-137  MatIsDiagonal                diagonal to a tolerance
//	NUM-LINALG-138  MatIsUpperTriangular         upper triangular to a tolerance
//	NUM-LINALG-139  MatIsLowerTriangular         lower triangular to a tolerance
//	NUM-LINALG-140  MatIsOrthogonal              QᵀQ = I to a tolerance
//	NUM-LINALG-141  MatIsPositiveDefinite        positive definite (Cholesky succeeds)
//	NUM-LINALG-142  MatIsPositiveSemidefinite    positive semidefinite (smallest eigenvalue ≥ −tol)
//	NUM-LINALG-143  MatIsSingular                numerical rank below n
//	NUM-LINALG-144  Mat2SolveCramer              2×2 system by Cramer's rule
//	NUM-LINALG-145  Mat3SolveCramer              3×3 system by Cramer's rule
//	NUM-LINALG-146  Mat2CharPoly                 2×2 characteristic polynomial
//	NUM-LINALG-147  Mat3CharPoly                 3×3 characteristic polynomial
//	NUM-LINALG-148  MatCharPolyFaddeevLeVerrier  n×n characteristic polynomial (Faddeev–LeVerrier)
//	NUM-LINALG-149  Mat2Eigen                    eigenvalues and eigenvectors of a general 2×2 (real or complex)
//	NUM-LINALG-150  MatSymmetricEigenJacobi      symmetric eigenvalues and vectors by cyclic Jacobi
//	NUM-LINALG-151  MatPowerIteration            dominant eigenpair by the power method
//	NUM-LINALG-152  Cholesky                     A = L·Lᵀ for symmetric positive definite A
//	NUM-LINALG-153  CholeskySolve                solve Ax = b from the Cholesky factor
//	NUM-LINALG-154  QRHouseholder                thin QR by Householder reflections
//	NUM-LINALG-155  GramSchmidtModified          thin QR by modified Gram–Schmidt
//	NUM-LINALG-156  LeastSquaresQR               min ‖Ax − b‖₂ by Householder QR
//	NUM-LINALG-157  ForwardSubstitution          solve Lx = b, L lower triangular
//	NUM-LINALG-158  BackSubstitution             solve Ux = b, U upper triangular
//	NUM-LINALG-159  TridiagonalSolve             tridiagonal system by the Thomas algorithm
//	NUM-LINALG-160  SVDJacobi                    thin SVD by one-sided Jacobi (Hestenes)
//	NUM-LINALG-161  MatPseudoInverse             Moore–Penrose pseudo-inverse from the SVD
//	NUM-LINALG-162  MatCondition2                2-norm condition number σ_max/σ_min
//	NUM-LINALG-163  MatNormFrobenius             Frobenius norm
//	NUM-LINALG-164  MatNorm1                     1-norm: largest absolute column sum
//	NUM-LINALG-165  MatNormInf                   ∞-norm: largest absolute row sum
//	NUM-LINALG-166  MatNorm2                     spectral norm σ_max
//	NUM-LINALG-167  CovarianceToCorrelation      correlation matrix and standard deviations from a covariance matrix
//	NUM-LINALG-168  CorrelationToCovariance      covariance matrix from correlations and standard deviations
//	NUM-LINALG-169  MatEigenvalues               all eigenvalues of a general real matrix (Hessenberg + Francis QR)
//	NUM-LINALG-170  Vandermonde                  Vandermonde matrix of ascending powers
//	NUM-LINALG-171  PolynomialFit                weighted least-squares polynomial of any degree (Vandermonde + QR)
//	NUM-LINALG-172  PolynomialRootsCompanion     polynomial roots as companion-matrix eigenvalues
//	NUM-LINALG-173  LeastSquaresMinNorm          minimum-norm least squares for any rank, by SVD
//	NUM-LINALG-174  MatLogDet                    sign and log|det A| from LU
//	NUM-LINALG-175  MatNormNuclear               nuclear (trace) norm Σσᵢ
//	NUM-LINALG-176  MatNormNeg1                  smallest absolute column sum
//	NUM-LINALG-177  MatNormNegInf                smallest absolute row sum
//	NUM-LINALG-178  MatNormNeg2                  smallest singular value
//	NUM-LINALG-179  VecNorm0                     count of nonzero entries
//	NUM-LINALG-180  MatCondition1                exact κ₁ = ‖A‖₁‖A⁻¹‖₁
//	NUM-LINALG-181  MatConditionInf              exact κ∞ = ‖A‖∞‖A⁻¹‖∞
//	NUM-LINALG-182  MatConditionFrobenius        κ_F = ‖A‖_F‖A⁻¹‖_F
//	NUM-LINALG-183  MatPowerInteger              Aᵏ for any integer k (negative via the inverse)
//	NUM-LINALG-184  MatMulChain                  product of a matrix chain in the optimal order
//	NUM-LINALG-185  VecMatMul                    row vector times matrix xᵀA
//	NUM-LINALG-186  MatDiagonal                  the offset-th diagonal of a matrix
//	NUM-LINALG-187  VecHadamard                  elementwise product of two vectors
//	NUM-LINALG-188  VecDivElementwise            elementwise quotient of two vectors
//	NUM-LINALG-189  SVDFull                      full SVD with square U and V
//	NUM-LINALG-190  PolynomialFitCovariance      polynomial fit with the coefficient covariance σ²(VᵀW²V)⁻¹ via R
//	NUM-LINALG-191  MatEigenvectors              eigenvalues and complex eigenvectors of a general real matrix (inverse iteration)
package linalg
