import { useFormContext, useFieldArray } from "react-hook-form";
import { motion } from "framer-motion";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import React from "react";
import { ScrollArea } from "../ui/scroll-area";

type SecurityQuestion = {
  question: string;
  answer: string;
};

type SecurityQuestionsFormData = {
  securityQuestions: SecurityQuestion[];
};

const securityQuestions = [
  "What was the name of your first pet?",
  "In what city were you born?",
  "What is your mother's maiden name?",
  "What high school did you attend?",
  "What was the make of your first car?",
  "What is your favorite movie?",
  "What is the name of your favorite childhood teacher?",
  "What is your favorite book?",
  "What is the name of the street you grew up on?",
  "What is your favorite food?",
];

export function SecurityQuestionsStep() {
  const {
    control,
    register,
    watch,
    formState: { errors },
    setValue,
    trigger,
  } = useFormContext<SecurityQuestionsFormData>();

  const { fields } = useFieldArray({
    control,
    name: "securityQuestions",
    rules: { minLength: 3, maxLength: 3 },
  });

  const watchSecurityQuestions = watch("securityQuestions");

  // Initialize exactly 3 questions once
  React.useEffect(() => {
    if (fields.length === 0) {
      const initialQuestions = Array(3).fill({ question: "", answer: "" });
      setValue("securityQuestions", initialQuestions);
    }
  }, [fields.length, setValue]);

  const getAvailableQuestions = (currentIndex: number) => {
    const selectedQuestions = watchSecurityQuestions?.reduce((acc: string[], q, idx) => {
      if (idx !== currentIndex && q.question) {
        acc.push(q.question);
      }
      return acc;
    }, []);
    
    return securityQuestions.filter(q => !selectedQuestions?.includes(q));
  };

  const handleQuestionChange = async (value: string, index: number) => {
    await setValue(`securityQuestions.${index}.question`, value, {
      shouldValidate: true,
      shouldDirty: true,
    });
    trigger(`securityQuestions.${index}.question`);
  };

  return (
    <div className="space-y-6 px-2">
      {fields.map((field, index) => {
        const availableQuestions = getAvailableQuestions(index);
        const currentValue = watchSecurityQuestions?.[index]?.question;

        return (
          <motion.div
            key={field.id}
            className="space-y-3"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
          >
            <Label className="text-lg font-medium text-gray-700">
              Security Question {index + 1}
            </Label>
            <Select
              value={currentValue || ""}
              onValueChange={(value) => handleQuestionChange(value, index)}
            >
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select a security question" />
              </SelectTrigger>
              <SelectContent>
                <ScrollArea className="h-[150px]">
                  {availableQuestions.map((question) => (
                    <SelectItem key={question} value={question}>
                      {question}
                    </SelectItem>
                  ))}
                  {currentValue && !availableQuestions.includes(currentValue) && (
                    <SelectItem key={currentValue} value={currentValue}>
                      {currentValue}
                    </SelectItem>
                  )}
                </ScrollArea>
              </SelectContent>
            </Select>
            {errors.securityQuestions?.[index]?.question && (
              <p className="text-red-500 text-sm mt-1">
                {errors.securityQuestions[index]?.question?.message}
              </p>
            )}
            <Input
              placeholder="Your answer"
              {...register(`securityQuestions.${index}.answer`)}
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
            />
            {errors.securityQuestions?.[index]?.answer && (
              <p className="text-red-500 text-sm mt-1">
                {errors.securityQuestions[index]?.answer?.message}
              </p>
            )}
          </motion.div>
        );
      })}
    </div>
  );
}
