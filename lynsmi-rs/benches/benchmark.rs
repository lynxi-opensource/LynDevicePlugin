use criterion::{criterion_group, criterion_main, Criterion};
use lynsmi::*;

fn get_device_props_benchmark(c: &mut Criterion) {
    let lib = Lib::try_default().unwrap();
    let symbols = Symbols::new(&lib).unwrap();

    let mut group = c.benchmark_group("lynsmi");
    group.sample_size(10);
    group.bench_function("get_device_props", |b| b.iter(|| symbols.get_props(0)));
    group.finish()
}

fn get_devices_benchmark(c: &mut Criterion) {
    let lib = Lib::try_default().unwrap();
    let smi = SMI::new(&lib).unwrap();

    let mut group = c.benchmark_group("lynsmi");
    group.sample_size(10);
    group.bench_function("get_devices", |b| {
        b.iter(|| {
            let results = &mut Vec::new();
            smi.get_devices(results);
        })
    });
    group.finish()
}

criterion_group!(benches, get_device_props_benchmark, get_devices_benchmark);
criterion_main!(benches);
