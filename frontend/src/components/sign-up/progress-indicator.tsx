"use client";

import { motion } from "framer-motion";

interface ProgressIndicatorProps {
  steps: string[];
  currentStep: number;
}

export function ProgressIndicator({
  steps,
  currentStep,
}: ProgressIndicatorProps) {
  return (
    <div className="relative py-8">
      {/* Steps labels */}
      <div className="flex justify-between mb-2">
        {steps.map((step, index) => (
          <motion.div
            key={`label-${step}`}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.2 }}
            className={`text-sm font-medium ${
              index <= currentStep ? "text-blue-600" : "text-gray-400"
            }`}
          >
            {step}
          </motion.div>
        ))}
      </div>

      {/* Progress bar background */}
      <div className="absolute top-[4.9rem] left-0 w-full h-[4px] bg-gray-200" />

      {/* Animated progress bar */}
      <motion.div
        className="absolute top-[4.9rem] left-0 h-[4px] bg-blue-600"
        initial={{ width: "0%" }}
        animate={{
          width: `${(currentStep / (steps.length - 1)) * 100}%`,
        }}
        transition={{ duration: 0.3, ease: "easeInOut" }}
      />

      {/* Step indicators */}
      <div className="relative flex justify-between">
        {steps.map((step, index) => (
          <div key={`step-${step}`} className="relative">
            {/* Completed step check mark or number */}
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{
                type: "spring",
                stiffness: 500,
                damping: 30,
                delay: index * 0.2,
              }}
              className={`relative z-10 flex items-center justify-center w-10 h-10 rounded-full border-2 ${
                index < currentStep
                  ? "bg-blue-600 border-blue-600"
                  : index === currentStep
                  ? "bg-white border-blue-600"
                  : "bg-white border-gray-200"
              }`}
            >
              {index < currentStep ? (
                <motion.svg
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  className="w-6 h-6 text-white"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 13l4 4L19 7"
                  />
                </motion.svg>
              ) : (
                <span
                  className={
                    index === currentStep ? "text-blue-600" : "text-gray-400"
                  }
                >
                  {index + 1}
                </span>
              )}
            </motion.div>

            {/* Pulse animation for current step */}
            {index === currentStep && (
              <motion.div
                className="absolute inset-0 rounded-full"
                initial={{ scale: 0.8, opacity: 0 }}
                animate={{
                  scale: [1, 1.3, 1],
                  opacity: [0.8, 0, 0],
                }}
                transition={{
                  duration: 2,
                  repeat: Infinity,
                  repeatType: "loop",
                }}
              >
                <div className="w-full h-full rounded-full bg-blue-600 opacity-20" />
              </motion.div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
