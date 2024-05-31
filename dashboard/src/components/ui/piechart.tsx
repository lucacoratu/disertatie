"use client"

import { PieChart } from "@mui/x-charts";
import { useTheme } from "next-themes";
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { FC } from "react";

type CustomPieChartProps = {
    labels: string[],
    values: number[],
    title: string,
}

const CustomPieChart: FC<CustomPieChartProps> = ({labels, values, title}): JSX.Element => {
    //Create the custom color pallete
    const pieChartColors = ['#003f5c', '#2f4b7c', '#665191', '#a05195', '#d45087', '#f95d6a', '#ff7c43', '#ffa600']

    const {theme, setTheme} = useTheme();

    const muiTheme = createTheme({
        palette: {
          mode: theme === "dark" ? "dark" : "light",
        },
    });


    return (
        <div className="py-6 flex flex-col gap-0 items-center border-2 b-2 rounded-xl bg-card">
            <h3 className="font-bold text-sm">
                {title}
            </h3>
            <ThemeProvider theme={muiTheme}>
                <PieChart series={[{
                    data: labels?.map((label, index) => (
                        {id: index, value: values[index], label: label}
                    )),
                    highlightScope: { faded: 'global', highlighted: 'item' },
                    faded: { innerRadius: 30, additionalRadius: -30, color: theme === "light" ? "#383838": "gray" },
                    },
                    ]}
                    width={350}
                    height={270}
                    colors={pieChartColors}
                    slotProps={{
                        legend: {
                            markGap: 4,
                            labelStyle: {
                                fontSize: 12,
                                fill: theme === "light" ? "black" : "whitesmoke",
                            },
                            direction: 'column',
                            position: { vertical: 'middle', horizontal: 'right' },
                            padding: 0,
                        }
                    }}
                />
            </ThemeProvider>
        </div>
    );
};

export default CustomPieChart;