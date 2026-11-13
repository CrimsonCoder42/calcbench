// Package signal transforms and filters one-dimensional sampled signals: the discrete Fourier transform
// (radix-2 Cooley–Tukey for powers of two, Bluestein's chirp-z for every other length, so every length is
// O(n log n)), the real-input transform and its inverse, frequency bins and spectrum shifts, the classical
// window functions, linear convolution and cross-correlation (direct, or through the FFT for long
// inputs), and the scalar Kalman filter. Two- and N-dimensional transforms are out of scope.
//
// ID short form: SIG (IDs are NUM-SIG-NNN).
//
// signal imports num/special (BesselI0 for the Kaiser window) and algo/reduce, nothing else from this
// module (build brief R1). Outputs go into a caller-supplied dst, grown when its capacity is short, as in
// num/linalg. Errors wrap the sentinels ErrEmpty, ErrNotFinite, ErrDimension, ErrOverflow, ErrOutOfRange
// and ErrMode; no function returns a NaN for finite input.
//
// Functions:
//
//	NUM-SIG-001  FFT        discrete Fourier transform of a complex sequence, any length
//	NUM-SIG-002  IFFT       inverse DFT (normalised by 1/n)
//	NUM-SIG-003  RFFT       DFT of a real sequence, bins 0…n/2
//	NUM-SIG-004  IRFFT      real sequence of length n from its n/2+1 RFFT bins
//	NUM-SIG-005  FFTFreq    frequency of each FFT bin
//	NUM-SIG-006  RFFTFreq   frequency of each RFFT bin
//	NUM-SIG-007  FFTShift   move the zero-frequency bin to the centre
//	NUM-SIG-008  IFFTShift  inverse of FFTShift
//	NUM-SIG-009  Hanning    Hann window
//	NUM-SIG-010  Hamming    Hamming window
//	NUM-SIG-011  Blackman   three-term Blackman window
//	NUM-SIG-012  Bartlett   Bartlett (triangular) window
//	NUM-SIG-013  Kaiser     Kaiser–Bessel window
//	NUM-SIG-014  Convolve   linear convolution, modes full, same, valid
//	NUM-SIG-015  Correlate  cross-correlation, modes full, same, valid
//	NUM-SIG-016  KalmanFilter1D  scalar Kalman filter (random walk + noise): Predict, Update, State, Variance, Gain, Reset
package signal
