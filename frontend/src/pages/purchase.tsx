import { PurchaseSection } from "@/components/license/purchase";
import { Snowflake } from "lucide-react";
import { FC, useEffect, useState } from "react";
import { motion } from "framer-motion";

const PurchasePage: FC = () => {
  const [isLoading, setIsLoading] = useState(true);
  useEffect(() => {
    const timer = setTimeout(() => {
      setIsLoading(false);
    }, 3000);
    return () => clearTimeout(timer);
  }, []);

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center h-screen bg-[#0078D4]">
        <motion.div
          initial={{ opacity: 0, scale: 0.5 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.5 }}
        >
          <Snowflake className="w-32 h-32 text-white" />
        </motion.div>
        <motion.div
          className="mt-8 flex justify-center items-center"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5, duration: 0.5 }}
        >
          <div className="w-12 h-12 rounded-full border-t-4 border-white animate-spin"></div>
        </motion.div>
        <motion.p
          className="mt-4 text-white text-xl"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1, duration: 0.5 }}
        >
          Getting things ready...
        </motion.p>
      </div>
    );
  }

  //
  return (
    <div className="min-h-screen bg-[#0078D4] flex-col text-white w-full flex items-center justify-center fixed">
      <div>
        <div className="flex items-center mb-8 flex-col space-y-2">
          <Snowflake className="w-12 h-12 mr-4" />
          <h1 className="text-3xl font-bold">Blizzflow Account Setup</h1>
        </div>
      </div>
      <PurchaseSection />
    </div>
  );
};


export default PurchasePage;