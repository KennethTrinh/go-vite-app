import { Toaster } from "@/components/ui/sonner";
import { queryClient } from "@/lib/query-client";
import { QueryClientProvider } from "@tanstack/react-query";
import { PropsWithChildren } from "react";
import { Outlet } from "react-router";
export const Component: React.FC<PropsWithChildren<{}>> = ({ children }) => {
  return (
    <QueryClientProvider client={queryClient}>
      <div className={`antialiased`}>
        {children}
        <Outlet />
        <Toaster />
      </div>
    </QueryClientProvider>
  );
};
