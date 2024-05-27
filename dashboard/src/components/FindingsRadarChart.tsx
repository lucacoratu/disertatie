"use client";

import { FC } from "react";
import {useTheme} from "next-themes";

import {
    Chart as ChartJS,
    RadialLinearScale,
    PointElement,
    LineElement,
    Filler,
    Tooltip,
    Legend,
} from 'chart.js';

import { Radar as RadarChart } from 'react-chartjs-2';

ChartJS.register(
    RadialLinearScale,
    PointElement,
    LineElement,
    Filler,
    Tooltip,
    Legend
);

ChartJS.defaults.borderColor = "rgba(245,245,245,0.9)";
  
let testdata = {
    labels: ['Thing 1', 'Thing 2', 'Thing 3'],
    datasets: [
        {
            label: '# of rule findings',
            data: [2, 9, 3],
            backgroundColor: 'rgba(109, 40, 217, 1)',
            borderColor: 'rgba(109, 40, 217, 0.9)',
            borderWidth: 1,
            showLabelBackdrop: false,
        },
    ],
};

const options = {
    scales: {
        r: {
           ticks: {
               display: false // Hides the labels in the middel (numbers)
           },
           pointLabels: {
            fontColor: '#443F5B',
            font: {
              size: 12
            }
          },
           grid: {
            color: ['whitesmoke']
           }
       }
    }
}

type FindingsRadarChartProps = {
    labels: string[],
    values: number[],
    title: string,
}

const FindingsRadarChart: FC<FindingsRadarChartProps> = ({labels, values, title}): JSX.Element => {
    testdata.labels = labels;
    testdata.datasets[0].data = values;
    testdata.datasets[0].label = title;

    const {theme, setTheme} = useTheme();
    if (theme === "light") {
        ChartJS.defaults.borderColor = "rgba(0,0,0,0.3)";
    }

    return (
        <div className="max-h-[343px] py-6 flex flex-col gap-3 items-center border-2 b-2 rounded-xl bg-card">
            <RadarChart data={testdata} options={options}/>
        </div>
    );
}

export default FindingsRadarChart;