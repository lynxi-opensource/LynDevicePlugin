use crate::errors::*;
use libloading::{Library, Symbol};
use serde::{Deserialize, Serialize};
use std::{
    ffi::{c_char, CStr, OsStr},
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

pub struct Symbols<'lib> {
    lib_get_device_cnt: Symbol<'lib, fn(&mut i32) -> i32>,
    lib_get_device_props_v1: Symbol<'lib, fn(i32, &mut lynDeviceProperties_t_v1) -> i32>,
    lib_get_device_props_v2: Symbol<'lib, fn(i32, &mut lynDeviceProperties_t_v2) -> i32>,
    driver_version: DriverVersion,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct Props {
    board: BoardProps,
    device: DeviceProps,
}

impl<'lib> Symbols<'lib> {
    pub fn new(lib: &'lib Lib) -> Result<Self> {
        unsafe {
            Ok(Self {
                lib_get_device_cnt: lib.0.get(b"lynGetDeviceCountSmi")?,
                lib_get_device_props_v1: lib.0.get(b"lynGetDeviceProperties")?,
                lib_get_device_props_v2: lib.0.get(b"lynGetDeviceProperties")?,
                driver_version: DriverVersion::local()?,
            })
        }
    }
    pub fn get_device_cnt(&self) -> Result<usize> {
        let mut cnt = 0;
        Error::check((self.lib_get_device_cnt)(&mut cnt))?;
        Ok(cnt as usize)
    }

    fn get_props_v1(&self, id: usize) -> Result<Props> {
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

    fn get_props_v2(&self, id: usize) -> Result<Props> {
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

    pub fn get_props(&self, id: usize) -> Result<Props> {
        match &self.driver_version {
            v if v < &V1_10_2 => self.get_props_v1(id),
            v if v >= &V1_10_2 => self.get_props_v2(id),
            _ => unreachable!(),
        }
    }
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

fn string_from_c(data: &[i8]) -> Result<String> {
    unsafe { Ok(CStr::from_ptr(data.as_ptr()).to_str()?.to_owned()) }
}

const V1_10_2: DriverVersion = DriverVersion(1, 10, 2);

#[derive(Debug, PartialEq, Eq, PartialOrd, Ord)]
struct DriverVersion(usize, usize, usize);

impl DriverVersion {
    pub fn local() -> Result<Self> {
        let output = Command::new("lynxi-smi").arg("-v").output()?;
        const PREFIX: &[u8; 13] = b"SMI version: ";
        let mut version_strs = output
            .stdout
            .strip_prefix(PREFIX)
            .ok_or(Error::StripVersionPrefix)?
            .split(|v| v == &b'.')
            .map(|v| String::from_utf8(v.to_owned()).map(|v| v.trim().parse::<usize>()));
        Ok(Self(
            version_strs.next().ok_or(Error::SplitVersion)???,
            version_strs.next().ok_or(Error::SplitVersion)???,
            version_strs.next().ok_or(Error::SplitVersion)???,
        ))
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_get_device_cnt() {
        let lib = Lib::try_default().unwrap();
        let symbols = Symbols::new(&lib).unwrap();
        println!("{:?}", symbols.get_device_cnt().unwrap())
    }

    #[test]
    fn test_get_device_prop() {
        let lib = Lib::try_default().unwrap();
        let symbols = Symbols::new(&lib).unwrap();
        println!("{:?}", symbols.get_props(0).unwrap())
    }
}