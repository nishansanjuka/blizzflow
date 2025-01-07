"use client";
import { FC, useRef, useState } from "react";
import { Button } from "../ui/button";
import { HandCoins, Loader2 } from "lucide-react";
import { motion } from "framer-motion";
import { Input } from "../ui/input";
import { ValidateLicense } from "@/blizzflow/backend/domain/services/license/licenseservice";
import { SaveLicense } from "@/blizzflow/backend/domain/handlers/license/licensehandler";
import { toast } from "sonner";
import { useAuth } from "@/hooks/use-auth";

export const PurchaseSection: FC = () => {
  const { setLicenseStatus } = useAuth();

  const licenseElement = useRef<HTMLInputElement | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handlePurchase = async () => {
    setIsLoading(true);
    try {
      if (licenseElement.current) {
        try {
          toast.loading("Validating license...");

          // Add artificial delay
          await new Promise((resolve) => setTimeout(resolve, 2000));

          await ValidateLicense(licenseElement.current.value);
          await SaveLicense(licenseElement.current.value);

          toast.dismiss();
          toast.success("License purchased successfully");
          setLicenseStatus(true);
        } catch (error) {
          toast.dismiss();
          toast.error("License validation failed");
        }
      }
    } catch (error) {
      if (typeof error === "string") {
        toast.error(error);
      }
    } finally {
      setIsLoading(false);
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
        <Button disabled={isLoading} onClick={handlePurchase}>
          {!isLoading ? (
            <>
              <HandCoins className="size-4" />
              <span>Purchase</span>
            </>
          ) : (
            <Loader2 className="animate-spin" />
          )}
        </Button>
      </div>
    </motion.div>
  );
};
