// Deklarasi minimal Web Bluetooth API — tidak termasuk lib.dom.d.ts bawaan TypeScript.
// Hanya mencakup subset yang dipakai untuk print langsung ke printer thermal BLE.
export {}

declare global {
  interface BluetoothRemoteGATTCharacteristic {
    uuid: string
    properties: {
      write: boolean
      writeWithoutResponse: boolean
    }
    writeValue(value: BufferSource): Promise<void>
    writeValueWithoutResponse(value: BufferSource): Promise<void>
  }

  interface BluetoothRemoteGATTService {
    uuid: string
    getCharacteristics(): Promise<BluetoothRemoteGATTCharacteristic[]>
  }

  interface BluetoothRemoteGATTServer {
    connected: boolean
    connect(): Promise<BluetoothRemoteGATTServer>
    disconnect(): void
    getPrimaryServices(): Promise<BluetoothRemoteGATTService[]>
  }

  interface BluetoothDevice {
    id: string
    name?: string
    gatt?: BluetoothRemoteGATTServer
  }

  interface BluetoothRequestDeviceFilter {
    services?: string[]
    name?: string
    namePrefix?: string
  }

  interface BluetoothRequestDeviceOptions {
    filters?: BluetoothRequestDeviceFilter[]
    acceptAllDevices?: boolean
    optionalServices?: string[]
  }

  interface Bluetooth {
    getDevices(): Promise<BluetoothDevice[]>
    requestDevice(options: BluetoothRequestDeviceOptions): Promise<BluetoothDevice>
  }

  interface Navigator {
    bluetooth?: Bluetooth
  }
}
