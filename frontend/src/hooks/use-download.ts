"use client";
import { useState, useCallback } from "react";

// Function to map MIME types to file extensions
const getFileExtensionFromMimeType = (mimeType: string): string => {
  const mimeToExtMap: Record<string, string> = {
    "audio/mpeg": "mp3",
    "audio/wav": "wav",
    "audio/ogg": "ogg",
    "image/jpeg": "jpg",
    "image/png": "png",
    "application/pdf": "pdf",
    // Add more mappings as needed
  };

  return mimeToExtMap[mimeType] || "bin"; // Default to .bin if not recognized
};

const generateRandomBase64 = (): string => {
  const array = new Uint8Array(3); // Generate 3 random bytes
  window.crypto.getRandomValues(array);
  return btoa(String.fromCharCode(...(array as any)));
};

const useDownload = () => {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isDownloaded, setIsDownloaded] = useState<boolean>(false);
  const [progress, setProgress] = useState<number>(0);

  const downloadFile = useCallback(
    async (url: string, name: string): Promise<void> => {
      setIsLoading(true);
      setIsDownloaded(false);
      setProgress(0);

      const timestamp = Date.now();
      const randomBase64 = generateRandomBase64();

      try {
        const response = await fetch(url);

        if (!response.ok) {
          throw new Error("Failed to download file");
        }

        const reader = response.body?.getReader();
        const contentLength = +(response.headers.get("Content-Length") || "0");
        let receivedLength = 0;
        const chunks: Uint8Array[] = [];

        while (true) {
          const { done, value } = await reader!.read();

          if (done) {
            break;
          }

          chunks.push(value);
          receivedLength += value.length;

          const percentComplete = Math.round(
            (receivedLength / contentLength) * 100
          );
          setProgress(percentComplete);
        }

        const blob = new Blob(chunks);
        const mimeType = response.headers.get("Content-Type") || "";
        const fileExtension = getFileExtensionFromMimeType(mimeType);
        const fileName = `${
          name.split(" ")[0]
        }_${timestamp}_${randomBase64}.${fileExtension}`;

        const downloadUrl = window.URL.createObjectURL(blob);
        const anchor = document.createElement("a");

        anchor.href = downloadUrl;
        anchor.download = fileName;
        anchor.click();

        window.URL.revokeObjectURL(downloadUrl);
        setIsDownloaded(true);
      } catch (error) {
        console.error("Download error:", error);
      } finally {
        setIsLoading(false);
      }
    },
    []
  );

  const reset = useCallback(() => {
    setIsLoading(false);
    setIsDownloaded(false);
    setProgress(0);
  }, []);

  return { isLoading, isDownloaded, progress, downloadFile, reset };
};

export default useDownload;
