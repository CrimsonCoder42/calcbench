package prob

import "testing"

func BenchmarkGumbelPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = GumbelPDF(1.5, 0, 1)
	}
}

func BenchmarkGumbelCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = GumbelCDF(1.5, 0, 1)
	}
}

func BenchmarkGumbelQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = GumbelQuantile(0.7, 0, 1)
	}
}

func BenchmarkGumbelMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = GumbelMeanVariance(0, 1)
	}
}

func BenchmarkLaplacePDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LaplacePDF(1.5, 0, 1)
	}
}

func BenchmarkLaplaceCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LaplaceCDF(1.5, 0, 1)
	}
}

func BenchmarkLaplaceQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LaplaceQuantile(0.7, 0, 1)
	}
}

func BenchmarkLaplaceMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = LaplaceMeanVariance(0, 1)
	}
}

func BenchmarkLogisticPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LogisticPDF(1.5, 0, 1)
	}
}

func BenchmarkLogisticCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LogisticCDF(1.5, 0, 1)
	}
}

func BenchmarkLogisticQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LogisticQuantile(0.7, 0, 1)
	}
}

func BenchmarkLogisticMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = LogisticMeanVariance(0, 1)
	}
}

func BenchmarkParetoPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ParetoPDF(2.5, 1, 3)
	}
}

func BenchmarkParetoCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ParetoCDF(2.5, 1, 3)
	}
}

func BenchmarkParetoQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ParetoQuantile(0.7, 1, 3)
	}
}

func BenchmarkParetoMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = ParetoMeanVariance(1, 3)
	}
}

func BenchmarkRayleighPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = RayleighPDF(1.5, 1)
	}
}

func BenchmarkRayleighCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = RayleighCDF(1.5, 1)
	}
}

func BenchmarkRayleighQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = RayleighQuantile(0.7, 1)
	}
}

func BenchmarkRayleighMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = RayleighMeanVariance(1)
	}
}

func BenchmarkTriangularPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = TriangularPDF(1.5, 0, 1, 2)
	}
}

func BenchmarkTriangularCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = TriangularCDF(1.5, 0, 1, 2)
	}
}

func BenchmarkTriangularQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = TriangularQuantile(0.7, 0, 1, 2)
	}
}

func BenchmarkTriangularMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = TriangularMeanVariance(0, 1, 2)
	}
}

func BenchmarkPowerFunctionPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PowerFunctionPDF(0.5, 3)
	}
}

func BenchmarkPowerFunctionCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PowerFunctionCDF(0.5, 3)
	}
}

func BenchmarkPowerFunctionQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PowerFunctionQuantile(0.7, 3)
	}
}

func BenchmarkPowerFunctionMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = PowerFunctionMeanVariance(3)
	}
}
