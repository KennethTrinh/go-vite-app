// adapted from https://github.com/hatchet-dev/hatchet/blob/main/frontend/app/src/pages/error/index.tsx
import { Button } from "@/components/ui/button";
import { Flag } from "lucide-react";
import { PropsWithChildren } from "react";
import { Link, useRouteError, useLocation } from "react-router";

const SUPPORT_URL = "https://example.com";

export default function ErrorBoundary() {
  const error = useRouteError();

  console.error(error);

  const Layout: React.FC<PropsWithChildren> = ({ children }) => (
    <div className="flex justify-center items-center flex-1 w-full h-[75vh]">
      <div className="flex flex-col space-y-2 text-center">{children}</div>
    </div>
  );

  if (
    error instanceof TypeError &&
    error.message.includes("Failed to fetch dynamically imported module:")
  ) {
    return (
      <Layout>
        <h1 className="text-2xl font-semibold tracking-tight">A New App Version is Available!</h1>
        <Button onClick={() => window.location.reload()}>Reload to Update</Button>
      </Layout>
    );
  }

  return <ErrorFallback />;
}

const ErrorFallback = () => {
  const { pathname } = useLocation();
  const returnUrl = pathname.includes("dashboard") ? "/dashboard" : "/";
  return (
    <div className="flex min-h-[70vh] flex-col items-center justify-center bg-background px-4 py-12 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-md text-center">
        <div className="mx-auto h-12 w-12 text-primary" />
        <h1 className="mt-4 text-5xl md:text-6xl font-bold tracking-tight text-foreground">
          Oops, something went wrong!
        </h1>
        <p className="mt-4 text-lg text-muted-foreground">
          We're sorry for the inconvenience, but an unexpected error occurred.
        </p>
        <div className="mt-6">
          <Link
            to={returnUrl}
            className="inline-flex items-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground shadow-sm transition-colors hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 mr-4"
          >
            Go to Homepage
          </Link>
          <Link
            to={SUPPORT_URL}
            target="_blank"
            className="inline-flex items-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground shadow-sm transition-colors hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 mr-4"
          >
            <Flag className="w-4 h-4 mr-2" />
            Report issue
          </Link>
        </div>
      </div>
    </div>
  );
};
