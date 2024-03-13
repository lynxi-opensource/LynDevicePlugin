use criterion::{criterion_group, criterion_main, Criterion};
use lynsmi::*;
use rayon::prelude::*;

fn get_device_props_benchmark(c: &mut Criterion) {
    let lib = Lib::try_default().unwrap();
    let symbols = PropsSymbols::new(&lib).unwrap();

    let mut group = c.benchmark_group("lynsmi");
    group.sample_size(10);
    group.bench_function("get_device_props", |b| b.iter(|| symbols.get_props(0)));
    group.finish()
}

fn get_devices_benchmark(c: &mut Criterion) {
    let lib = Lib::try_default().unwrap();
    let smi_common = CommonSymbols::new(&lib).unwrap();
    let smi = PropsSymbols::new(&lib).unwrap();
    let cnt = smi_common.get_device_cnt().unwrap();

    let mut group = c.benchmark_group("lynsmi");
    group.sample_size(10);
    group.bench_function("get_devices", |b| {
        b.iter(|| {
            (0..cnt).into_par_iter().for_each(|id| {
                smi.get_props(id).unwrap();
            });
        })
    });
    group.finish()
}

criterion_group!(benches, get_device_props_benchmark, get_devices_benchmark);
criterion_main!(benches);
