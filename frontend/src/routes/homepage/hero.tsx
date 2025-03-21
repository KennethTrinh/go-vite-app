import { useEffect, useMemo, useState } from "react";
import { AlarmCheck, MoveRight, PhoneCall } from "lucide-react";
import { toast } from "sonner";
import { useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { listItems } from "@/api/items";

// creds: https://21st.dev/serafim/animated-hero/default
export const Hero = () => {
  const [titleNumber, setTitleNumber] = useState(0);
  const titles = useMemo(() => ["amazing", "new", "wonderful", "beautiful", "smart"], []);
  const navigate = useNavigate();

  useEffect(() => {
    const timeoutId = setTimeout(() => {
      if (titleNumber === titles.length - 1) {
        setTitleNumber(0);
      } else {
        setTitleNumber(titleNumber + 1);
      }
    }, 2000);
    return () => clearTimeout(timeoutId);
  }, [titleNumber, titles]);

  return (
    <div className="w-full">
      <div className="container mx-auto">
        <div className="flex gap-8 py-20 lg:py-40 items-center justify-center flex-col">
          <div>
            <Button variant="secondary" size="sm" className="gap-4">
              Read our launch article <MoveRight className="w-4 h-4" />
            </Button>
          </div>
          <div className="flex gap-4 flex-col">
            <h1 className="text-5xl md:text-7xl max-w-2xl tracking-tighter text-center font-regular">
              <span className="text-spektr-cyan-50">This is something</span>
              <span className="relative flex w-full justify-center overflow-hidden text-center md:pb-4 md:pt-1">
                &nbsp;
                <span className="relative flex w-full justify-center overflow-hidden text-center md:pb-4 md:pt-1">
                  &nbsp;
                  {titles.map((title, index) => (
                    <span
                      key={index}
                      className={`absolute font-semibold transition-all duration-500 ${
                        titleNumber === index
                          ? "opacity-100 translate-y-0"
                          : titleNumber > index
                          ? "opacity-0 -translate-y-[150px]"
                          : "opacity-0 translate-y-[150px]"
                      }`}
                    >
                      {title}
                    </span>
                  ))}
                </span>
              </span>
            </h1>

            <p className="text-lg md:text-xl leading-relaxed tracking-tight text-muted-foreground max-w-2xl text-center">
              Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
              incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud
              exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure
              dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
              Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt
              mollit anim id est laborum.
            </p>
          </div>
          <div className="flex flex-row gap-3">
            <Button
              size="lg"
              className="gap-4"
              variant="destructive"
              onClick={async () => {
                const response = await listItems();
                if ("error" in response) {
                  toast.error(response.error.message);
                  return;
                }
                toast.success(`Fetched ${response?.items?.length} items`);
              }}
            >
              Try 429s (Prod only) <AlarmCheck className="w-4 h-4" />
            </Button>

            <Button
              size="lg"
              className="gap-4"
              variant="outline"
              onClick={() => toast.success("Hello, this is a success toast message")}
            >
              Try toast <PhoneCall className="w-4 h-4" />
            </Button>
            <Button
              size="lg"
              className="gap-4"
              onClick={() => {
                navigate("/dynamic");
              }}
            >
              See API Call <MoveRight className="w-4 h-4" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};
