// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "@wailsio/runtime";

export function ReadLicense(): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2216710016) as any;
    return $resultPromise;
}

export function SaveLicense(licenseKey: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2402892247, licenseKey) as any;
    return $resultPromise;
}
