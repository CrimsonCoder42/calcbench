# Kahan

**Citation**: Kahan, "How Futile are Mindless Assessments of Roundoff in Floating-Point Computation?", University of California, Berkeley, 11 January 2006, https://people.eecs.berkeley.edu/~wkahan/Mindless.pdf.

**Taken**: §12 "Mangled Angles", pp. 46–47: the angle between vectors ∠(x, y) = 2 arctan(‖x‖y‖ − ‖x‖y‖ / ‖x‖y‖ + ‖x‖y‖); the arccos formula loses half the digits near 0 and π, the cross-product arcsin formula near π/2.

**Used by**: NUM-LINALG-009, NUM-LINALG-020, NUM-LINALG-021, NUM-LINALG-047.

**Checked**: PDF text extracted and read 2026-11-12. Tests: nearly parallel vectors give atan(1e−10) to 1e−25 (Vec2AngleBetween_nearlyParallel), the n-D angle agrees with the 3D one.
