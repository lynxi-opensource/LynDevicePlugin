pub mod errors;
pub mod library;
pub use library::*;

use rayon::{prelude::*, ThreadPool};

pub struct SMI<'a> {
    symbols: Symbols<'a>,
    device_cnt: usize,
    thread_pool: ThreadPool,
}

impl<'a> SMI<'a> {
    pub fn new(lib: &'a Lib) -> Result<Self, errors::Error> {
        let symbols = Symbols::new(lib)?;
        let device_cnt = symbols.get_device_cnt()?;

        let thread_pool = rayon::ThreadPoolBuilder::new()
            .num_threads(device_cnt)
            .build()?;

        Ok(Self {
            symbols,
            device_cnt,
            thread_pool,
        })
    }

    pub fn get_devices(&self, results: &mut Vec<Result<Props, errors::Error>>) {
        self.thread_pool.install(|| {
            (0..self.device_cnt)
                .into_par_iter()
                .map(|id| self.symbols.get_props(id))
                .collect_into_vec(results);
        });
    }
}
