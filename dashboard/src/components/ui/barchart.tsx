"use client"

import { BarChart } from "@mui/x-charts";
import { axisClasses } from '@mui/x-charts/ChartsAxis';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { useTheme } from "next-themes";
import { FC } from "react";

type CustomBarChartProps = {
    labels: string[],
    values: number[],
    title: string,
}

const CustomPieChart: FC<CustomBarChartProps> = ({labels, values, title}): JSX.Element => {
    //Create the custom color pallete
    const pieChartColors = ['#003f5c', '#2f4b7c', '#665191', '#a05195', '#d45087', '#f95d6a', '#ff7c43', '#ffa600']

    const {theme, setTheme} = useTheme();

    const muiTheme = createTheme({
        palette: {
          mode: theme === "dark" ? "dark" : "light",
        },
    });

    return (
        <div className="py-6 px-3 flex flex-col gap-0 items-center border-2 rounded-xl bg-card">
            <h3 className="font-bold text-sm">
                {title}
            </h3>
            {values && values.length > 0 &&
            <ThemeProvider theme={muiTheme}>
            <BarChart
                series={[
                    {
                        data: values,
                        color: pieChartColors[1],
                    }
                ]}
                
                xAxis={[
                    {
                        scaleType: "band",
                        data: labels,
                        labelStyle: {
                            fontSize: 14,
                            // transform: `translateY(${
                            //     // Hack that should be added in the lib latter.
                            //     5 * Math.abs(Math.sin((Math.PI * 45) / 180))
                            // }px)`,
                            fill: theme === "light" ? "black": "whitesmoke"
                        },
                        tickLabelStyle: {
                            angle: 15,
                            textAnchor: 'middle',
                            fontSize: 9,
                            fill: theme === "light" ? "black": "whitesmoke"
                        },
                    },
                ]}
                yAxis={[{
                    label: "count",
                    labelStyle: {
                        fill: theme === "light" ? "black": "whitesmoke",
                    }
                }]}
                //width={400}
                height={270}
                tooltip={{ trigger: 'axis' }}

                sx={() => ({
                    [`.${axisClasses.root}`]: {
                      [`.${axisClasses.tick}, .${axisClasses.line}`]: {
                        stroke: theme === "light" ? "black": "whitesmoke",
                        strokeWidth: 1,
                      },
                      [`.${axisClasses.tickLabel}`]: {
                        fill: theme === "light" ? "black": "whitesmoke",
                      },
                      [`.${axisClasses.label}`] : {
                        fill: theme === "light" ? "black": "whitesmoke",
                      }
                    },
                })}
            />
            </ThemeProvider>
            }
        </div>
    );
};

export default CustomPieChart;