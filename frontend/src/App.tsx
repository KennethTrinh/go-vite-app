import { createBrowserRouter, createRoutesFromElements, Route, RouterProvider } from "react-router";
import { Component as Root } from "@/routes/root-layout";
import ErrorBoundary from "@/components/error-boundary";

export default function App() {
  return <RouterProvider router={router} />;
}

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route
      path="/"
      element={<Root />}
      errorElement={
        <Root>
          <ErrorBoundary />{" "}
        </Root>
      }
      hydrateFallbackElement={<></>}
    >
      <Route index lazy={() => import("@/routes/homepage/homepage")} />
      <Route path="dynamic" lazy={() => import("@/routes/dynamic/dynamic")} />
    </Route>
  )
);
