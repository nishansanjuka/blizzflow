"use client";
import { FC, useRef } from "react";
import { Button } from "../ui/button";
import { HandCoins } from "lucide-react";
import { motion } from "framer-motion";
import { Input } from "../ui/input";
import { ValidateLicense } from "@/blizzflow/backend/domain/services/license/licenseservice";
import { SaveLicense } from "@/blizzflow/backend/domain/handlers/license/licensehandler";

export const PurchaseSection: FC = () => {
  const licenseElement = useRef<HTMLInputElement | null>(null);
  const handlePurchase = async () => {
    try {
      if (licenseElement.current) {
        const res = await ValidateLicense(licenseElement.current.value);
        console.log(res);
        await SaveLicense(licenseElement.current.value);
      }
    } catch (error: any) {
      console.log(error);
    }
  };
  return (
    <motion.div
      initial={{ y: 20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.3, ease: "easeOut" }}
      className="w-1/2"
    >
      <div className="w-full flex flex-col space-y-2">
        <Input
          ref={licenseElement}
          className="border-blue-500 placeholder:text-blue-300 focus:ring-blue-400 focus:border-blue-400 text-center focus-visible:ring-blue-400 w-full focus:outline-none"
          placeholder="XXXX-XXXX-XXXX-XXXX"
        />
        <Button onClick={handlePurchase}>
          <HandCoins className="size-4" />
          <span>Purchase</span>
        </Button>
      </div>
    </motion.div>
  );
};
