use crate::errors::*;
use libloading::{Library, Symbol};
use regex::Regex;
use serde::{Deserialize, Serialize};
use std::{
    ffi::{c_char, c_int, CStr, OsStr},
    fmt::Debug,
    mem::zeroed,
    process::Command,
};

const ARRAY_MAX_LEN: usize = 40;
const PROCESS_NAME_LEN: usize = 64;
const PROCESS_COUNT_LIMIT: usize = 16;

#[repr(C)]
#[allow(non_snake_case)]
struct lynDeviceProperties_t_v1 {
    boardProductName: [c_char; ARRAY_MAX_LEN],
    boardBrand: [c_char; ARRAY_MAX_LEN],
    boardFirmwareVersion: [c_char; ARRAY_MAX_LEN],
    boardProductNumber: [c_char; ARRAY_MAX_LEN],
    boardSerialNumber: [c_char; ARRAY_MAX_LEN],
    boardId: u32,
    boardChipCount: u32,
    boardPowerDraw: f32,
    boardPowerLimit: f32,
    boardVoltage: f32,

    deviceName: [c_char; ARRAY_MAX_LEN],
    deviceUuid: [c_char; ARRAY_MAX_LEN],
    deviceApuClockFrequency: u64,
    deviceApuClockFrequencyLimit: u64,
    deviceArmClockFrequency: u64,
    deviceArmClockFrequencyLimit: u64,
    deviceMemClockFrequency: u64,
    deviceMemClockFrequencyLimit: u64,
    deviceMemoryUsed: u64,
    deviceMemoryTotal: u64,
    deviceTemperatureCurrent: i32,
    deviceTemperatureSlowdown: i32,
    deviceTemperatureLimit: i32,
    deviceApuUsageRate: u32,
    deviceArmUsageRate: u32,
    deviceVicUsageRate: u32,
    deviceIpeUsageRate: u32,
    deviceEccStat: u32,
    deviceDdrErrorCount: u32,
    deviceDdrNoErrorCount: u32,
    deviceVoltage: f32,

    processCount: u32,
    pid: [u32; PROCESS_COUNT_LIMIT],
    processName: [[u8; PROCESS_NAME_LEN]; PROCESS_COUNT_LIMIT],
    processUseMemory: [u64; PROCESS_COUNT_LIMIT],
}

#[repr(C)]
#[allow(non_snake_case)]
struct lynDeviceProperties_t_v2 {
    boardProductName: [c_char; ARRAY_MAX_LEN],
    boardBrand: [c_char; ARRAY_MAX_LEN],
    boardFirmwareVersion: [c_char; ARRAY_MAX_LEN],
    boardSerialNumber: [c_char; ARRAY_MAX_LEN],
    boardId: u32,
    boardChipCount: u32,
    boardPowerDraw: f32,
    boardPowerLimit: f32,
    boardVoltage: f32,

    deviceName: [c_char; ARRAY_MAX_LEN],
    deviceUuid: [c_char; ARRAY_MAX_LEN],
    deviceApuClockFrequency: u64,
    deviceApuClockFrequencyLimit: u64,
    deviceArmClockFrequency: u64,
    deviceArmClockFrequencyLimit: u64,
    deviceMemClockFrequency: u64,
    deviceMemClockFrequencyLimit: u64,
    deviceMemoryUsed: u64,
    deviceMemoryTotal: u64,
    deviceTemperatureCurrent: i32,
    deviceTemperatureSlowdown: i32,
    deviceTemperatureLimit: i32,
    deviceApuUsageRate: u32,
    deviceArmUsageRate: u32,
    deviceVicUsageRate: u32,
    deviceIpeUsageRate: u32,
    deviceEccStat: u32,
    deviceDdrErrorCount: u32,
    deviceDdrNoErrorCount: u32,
    deviceVoltage: f32,

    processCount: u32,
    pid: [u32; PROCESS_COUNT_LIMIT],
    processName: [[u8; PROCESS_NAME_LEN]; PROCESS_COUNT_LIMIT],
    processUseMemory: [u64; PROCESS_COUNT_LIMIT],
}

pub struct Lib(Library);

impl Lib {
    const DEFAULT_FILENAME: &str = "/usr/lib/libLYNSMICLIENTCOMM.so";

    pub fn new<P>(filename: P) -> Result<Self>
    where
        P: AsRef<OsStr>,
    {
        unsafe { Ok(Lib(Library::new(filename)?)) }
    }

    pub fn try_default() -> Result<Self> {
        Self::new(Self::DEFAULT_FILENAME)
    }
}

#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
pub enum P2PMode {
    NonSupport,
    P2PLinkPIX,
    P2PLinkPXB,
    P2PLinkPHB,
    P2PLinkSYS,
}

#[derive(Clone)]
pub struct CommonSymbols<'lib> {
    lib_get_device_cnt: Symbol<'lib, fn(&mut i32) -> c_int>,
}

impl<'lib> CommonSymbols<'lib> {
    pub fn new(lib: &'lib Lib) -> Result<Self> {
        unsafe {
            Ok(Self {
                lib_get_device_cnt: lib.0.get(b"lynGetDeviceCountSmi")?,
            })
        }
    }

    pub fn get_device_cnt(&self) -> Result<i32> {
        let mut cnt = 0;
        Error::check((self.lib_get_device_cnt)(&mut cnt))?;
        Ok(cnt)
    }
}

#[derive(Clone)]
pub struct TopologySymbols<'lib> {
    lib_device_show_topology: Symbol<'lib, fn() -> c_int>,
    lib_get_device_p2p_attr: Symbol<'lib, fn(c_int, c_int, &mut P2PMode, &mut c_int) -> c_int>,
}

impl<'lib> TopologySymbols<'lib> {
    pub fn new(lib: &'lib Lib) -> Result<Self> {
        unsafe {
            Ok(Self {
                lib_device_show_topology: lib.0.get(b"lynDeviceShowTopologyS")?,
                lib_get_device_p2p_attr: lib.0.get(b"lynGetDeviceP2PAttrS")?,
            })
        }
    }

    pub fn show_topology(&self) -> Result<()> {
        let f = &self.lib_device_show_topology;
        Error::check(f())
    }

    pub fn get_device_p2p_attr(&self, src_device: i32, dst_device: i32) -> Result<P2PAttr> {
        let f = &self.lib_get_device_p2p_attr;
        let mut mode = P2PMode::NonSupport;
        let mut dist: c_int = 0;
        Error::check(f(src_device, dst_device, &mut mode, &mut dist))?;
        Ok(P2PAttr { mode, dist })
    }
}

#[derive(Clone)]
pub struct PropsSymbols<'lib> {
    lib_get_device_props_v1: Symbol<'lib, fn(i32, &mut lynDeviceProperties_t_v1) -> c_int>,
    lib_get_device_props_v2: Symbol<'lib, fn(i32, &mut lynDeviceProperties_t_v2) -> c_int>,
    driver_version: DriverVersion,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct Props {
    pub board: BoardProps,
    pub device: DeviceProps,
}

impl<'lib> PropsSymbols<'lib> {
    pub fn new(lib: &'lib Lib) -> Result<Self> {
        unsafe {
            Ok(Self {
                lib_get_device_props_v1: lib.0.get(b"lynGetDeviceProperties")?,
                lib_get_device_props_v2: lib.0.get(b"lynGetDeviceProperties")?,
                driver_version: DriverVersion::local()?,
            })
        }
    }

    pub fn get_driver_version(&self) -> &DriverVersion {
        &self.driver_version
    }

    fn get_props_v1(&self, id: i32) -> Result<Props> {
        let mut c_device_prop: lynDeviceProperties_t_v1 = unsafe { zeroed() };
        Error::check((self.lib_get_device_props_v1)(
            id as i32,
            &mut c_device_prop,
        ))?;
        Ok(Props {
            board: BoardProps {
                product_name: string_from_c(c_device_prop.boardProductName.as_ref())?,
                brand: string_from_c(c_device_prop.boardBrand.as_ref())?,
                serial_number: string_from_c(c_device_prop.boardSerialNumber.as_ref())?,
                id: c_device_prop.boardId,
                chip_count: c_device_prop.boardChipCount,
                power_draw: c_device_prop.boardPowerDraw,
            },
            device: DeviceProps {
                name: string_from_c(c_device_prop.deviceName.as_ref())?,
                uuid: string_from_c(c_device_prop.deviceUuid.as_ref())?,
                memory_used: c_device_prop.deviceMemoryUsed,
                memory_total: c_device_prop.deviceMemoryTotal,
                temperature: c_device_prop.deviceTemperatureCurrent,
                apu_usage: c_device_prop.deviceApuUsageRate,
                arm_usage: c_device_prop.deviceArmUsageRate,
                vic_usage: c_device_prop.deviceVicUsageRate,
                ipe_usage: c_device_prop.deviceIpeUsageRate,
            },
        })
    }

    fn get_props_v2(&self, id: i32) -> Result<Props> {
        let mut c_device_prop: lynDeviceProperties_t_v2 = unsafe { zeroed() };
        Error::check((self.lib_get_device_props_v2)(
            id as i32,
            &mut c_device_prop,
        ))?;
        Ok(Props {
            board: BoardProps {
                product_name: string_from_c(c_device_prop.boardProductName.as_ref())?,
                brand: string_from_c(c_device_prop.boardBrand.as_ref())?,
                serial_number: string_from_c(c_device_prop.boardSerialNumber.as_ref())?,
                id: c_device_prop.boardId,
                chip_count: c_device_prop.boardChipCount,
                power_draw: c_device_prop.boardPowerDraw,
            },
            device: DeviceProps {
                name: string_from_c(c_device_prop.deviceName.as_ref())?,
                uuid: string_from_c(c_device_prop.deviceUuid.as_ref())?,
                memory_used: c_device_prop.deviceMemoryUsed,
                memory_total: c_device_prop.deviceMemoryTotal,
                temperature: c_device_prop.deviceTemperatureCurrent,
                apu_usage: c_device_prop.deviceApuUsageRate,
                arm_usage: c_device_prop.deviceArmUsageRate,
                vic_usage: c_device_prop.deviceVicUsageRate,
                ipe_usage: c_device_prop.deviceIpeUsageRate,
            },
        })
    }

    pub fn get_props(&self, id: i32) -> Result<Props> {
        match &self.driver_version {
            v if v < &V1_10_2 => self.get_props_v1(id),
            v if v >= &V1_10_2 => self.get_props_v2(id),
            _ => unreachable!(),
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
pub struct P2PAttr {
    pub mode: P2PMode,
    pub dist: i32,
}

#[non_exhaustive]
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct BoardProps {
    pub product_name: String,
    pub brand: String,
    pub serial_number: String,
    pub id: u32,
    pub chip_count: u32,
    pub power_draw: f32,
}

#[non_exhaustive]
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct DeviceProps {
    pub name: String,
    pub uuid: String,
    pub memory_used: u64,
    pub memory_total: u64,
    pub temperature: i32,
    pub apu_usage: u32,
    pub arm_usage: u32,
    pub vic_usage: u32,
    pub ipe_usage: u32,
}

fn string_from_c(data: &[c_char]) -> Result<String> {
    unsafe { Ok(CStr::from_ptr(data.as_ptr()).to_str()?.to_owned()) }
}

const V1_10_2: DriverVersion = DriverVersion(1, 10, 2);

#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
pub struct DriverVersion(pub usize, pub usize, pub usize);

impl DriverVersion {
    pub fn local() -> Result<Self> {
        let re = Regex::new(r"SMI version: (\d+)\.(\d+)\.(\d+)").unwrap();
        let output = Command::new("lynxi-smi").arg("-v").output()?;
        let output = String::from_utf8(output.stdout)?;
        let captures = re
            .captures(&output)
            .ok_or_else(|| Error::NoVersionInfo(output.clone()))?;
        Ok(Self(
            captures[1].parse::<usize>()?,
            captures[2].parse::<usize>()?,
            captures[3].parse::<usize>()?,
        ))
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_get_device_cnt() {
        let lib = Lib::try_default().unwrap();
        let symbols = CommonSymbols::new(&lib).unwrap();
        println!("{:?}", symbols.get_device_cnt().unwrap())
    }

    #[test]
    fn test_get_device_prop() {
        let lib = Lib::try_default().unwrap();
        let symbols = PropsSymbols::new(&lib).unwrap();
        println!("{:?}", symbols.get_props(0).unwrap())
    }

    #[test]
    fn test_show_topology() {
        let lib = Lib::try_default().unwrap();
        let symbols = TopologySymbols::new(&lib).unwrap();
        println!("{:?}", symbols.show_topology());
    }

    #[test]
    fn test_get_device_p2p_attr() {
        let lib = Lib::try_default().unwrap();
        let common_symbols = CommonSymbols::new(&lib).unwrap();
        let symbols = TopologySymbols::new(&lib).unwrap();
        let device_cnt = common_symbols.get_device_cnt().unwrap();
        for src_device in 0..device_cnt {
            for dst_device in (src_device + 1)..device_cnt {
                println!(
                    "{} {} {:?}",
                    src_device,
                    dst_device,
                    symbols
                        .get_device_p2p_attr(src_device as i32, dst_device as i32)
                        .unwrap()
                );
            }
        }
    }

    #[test]
    fn test_get_driver_version() {
        let version = DriverVersion::local().unwrap();
        println!("{:?}", version);
    }
}
