// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as model$0 from "../../model/models.js";

export function CreateUser(username: string, password: string): Promise<model$0.User | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2101710880, username, password) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function DeleteUser(userID: number): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1059264179, userID) as any;
    return $resultPromise;
}

export function GetUserByID(userID: number): Promise<model$0.User | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(215609246, userID) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function GetUserByUsername(username: string): Promise<model$0.User | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3485513747, username) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function UpdateUser(user: model$0.User | null): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(618010709, user) as any;
    return $resultPromise;
}

// Private type creation functions
const $$createType0 = model$0.User.createFrom;
const $$createType1 = $Create.Nullable($$createType0);