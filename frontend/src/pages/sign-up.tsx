"use client";

import { useState } from "react";
import { useForm, FormProvider } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { motion, AnimatePresence } from "framer-motion";
import { ProgressIndicator } from "@/components/sign-up/progress-indicator";
import { UsernamePasswordStep } from "@/components/sign-up/username-password-section";
import { SecurityQuestionsStep } from "@/components/sign-up/security-questions-step";
import { Loader2, Snowflake } from "lucide-react";
import { toast } from "sonner";
import { UserService } from "@/blizzflow/backend/domain/services/user";
import { useAuth } from "@/hooks/use-auth";
import { Window } from "@wailsio/runtime";
import { Navigate, useNavigate } from "react-router-dom";

const schema = z.object({
  username: z.string().min(3, "Username must be at least 3 characters"),
  password: z.string().min(8, "Password must be at least 8 characters"),
  securityQuestions: z
    .array(
      z.object({
        question: z.string(),
        answer: z.string().min(1, "Answer is required"),
      })
    )
    .length(3, "Please answer all 3 security questions"),
});

type FormData = z.infer<typeof schema>;

const steps = ["Account Creation", "Security Questions"];

export default function ProfileSetup() {
  const [currentStep, setCurrentStep] = useState(0);
  const [loading, setLoading] = useState(false);
  const { setSecurityQuestions } = useAuth();
  const methods = useForm<FormData>({
    resolver: zodResolver(schema),
    mode: "onChange",
  });

  const navigate = useNavigate();

  const username = methods.watch("username");
  const password = methods.watch("password");

  const handleNextStep = () => {
    if (currentStep < steps.length - 1) {
      if (currentStep === 0 && username && password) {
        nextStep();
      } else {
        methods.trigger();
      }
    }
  };

  const onSubmit = async (data: FormData) => {
    setLoading(true);
    try {
      toast.loading("Signing in...");
      // Add artificial delay
      await new Promise((resolve) => setTimeout(resolve, 2000));
      const SecurityQuestionsRecord = data.securityQuestions.reduce(
        (acc, { question, answer }) => ({
          ...acc,
          [question]: answer,
        }),
        {} as Record<string, string>
      );
      await UserService.CreateUser(data.username, data.password);
      localStorage.setItem("username", data.username);
      await setSecurityQuestions(data.username, SecurityQuestionsRecord);
      toast.success("Account created successfully");
      navigate("/callback", {
        viewTransition: true,
      });
    } catch (error) {
      if (typeof error === "string") {
        toast.error(error.split(":")[1]);
      }
    } finally {
      toast.dismiss();
      setLoading(false);
    }
  };

  const nextStep = () =>
    setCurrentStep((prev) => Math.min(prev + 1, steps.length - 1));
  const prevStep = () => setCurrentStep((prev) => Math.max(prev - 1, 0));

  return (
    <div className="min-h-screen to-white flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="rounded-lg w-full max-w-2xl p-8"
      >
        <h1 className="text-3xl font-semibold mb-6 text-gray-800 flex items-center space-x-2">
          <Snowflake className="size-8 text-blue-500" />
          <span>Account Setup</span>
        </h1>
        <ProgressIndicator steps={steps} currentStep={currentStep} />
        <FormProvider {...methods}>
          <form onSubmit={methods.handleSubmit(onSubmit)} className="space-y-6">
            <AnimatePresence mode="wait">
              {currentStep === 0 && (
                <motion.div
                  key="username-password"
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: 20 }}
                  transition={{ duration: 0.3 }}
                >
                  <UsernamePasswordStep />
                </motion.div>
              )}
              {currentStep === 1 && (
                <motion.div
                  key="security-questions"
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: 20 }}
                  transition={{ duration: 0.3 }}
                >
                  <SecurityQuestionsStep />
                </motion.div>
              )}
            </AnimatePresence>
            <div className="flex justify-between mt-8">
              {currentStep > 0 && (
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  type="button"
                  onClick={prevStep}
                  className="px-6 py-2 bg-gray-200 text-gray-800 rounded-md hover:bg-gray-300 transition-colors"
                >
                  Back
                </motion.button>
              )}
              <motion.button
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                type={currentStep < steps.length - 1 ? "button" : "submit"}
                onClick={handleNextStep}
                disabled={loading}
                className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors ml-auto"
              >
                {!loading ? (
                  <>{currentStep < steps.length - 1 ? "Next" : "Finish"}</>
                ) : (
                  <Loader2 className="animate-spin size-4" />
                )}
              </motion.button>
            </div>
          </form>
        </FormProvider>
      </motion.div>
    </div>
  );
}
