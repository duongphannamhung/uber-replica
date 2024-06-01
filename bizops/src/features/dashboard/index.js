import DashboardStats from './components/DashboardStats'
import AmountStats from './components/AmountStats'
import PageStats from './components/PageStats'

import UserGroupIcon  from '@heroicons/react/24/outline/UserGroupIcon'
import UsersIcon  from '@heroicons/react/24/outline/UsersIcon'
import CircleStackIcon  from '@heroicons/react/24/outline/CircleStackIcon'
import CreditCardIcon  from '@heroicons/react/24/outline/CreditCardIcon'
import UserChannels from './components/UserChannels'
import LineChart from './components/LineChart'
import BarChartCustom from './components/BarChartCustom'
import DashboardTopBar from './components/DashboardTopBar'
import { useDispatch } from 'react-redux'
import {showNotification} from '../common/headerSlice'
import DoughnutChart from './components/DoughnutChart'
import { useState, useEffect } from 'react'
import axios from 'axios';

const today = new Date();
const yesterday = new Date(today);
yesterday.setDate(yesterday.getDate() - 1);

let defaultStartDate = yesterday.toISOString().split('T')[0];
let defaultEndDate = today.toISOString().split('T')[0];

function toVND(value) {
    if (!value) return value;
    console.debug('value to VND: ', value)
    return value.toLocaleString('vi', {style : 'currency', currency : 'VND'});
}

function Dashboard(){
    const [statsData, setStatsData] = useState([
        {title : "Active Users", value : Math.round(Math.random() * 1000), icon : <UsersIcon className='w-8 h-8'/>, description : `↙ ${Math.round(Math.random() * 100)} (${Math.round(Math.random() * 10)}%)`}
    ])

    const fetchNewUsers = (startDate=defaultStartDate, endDate=defaultEndDate) => {
        axios.post('/api/dashboard/calculate-new-users', {
            start_date: startDate,
            end_date: endDate
        })
        .then(response => {
            const data = response.data;
            const arrow = data.is_increase ? '↗︎' : '↙';
            const newUserStat = {
                title : "New Users", 
                value : data.value, 
                icon : <UserGroupIcon className='w-8 h-8'/>, 
                description : `${arrow} ${data.difference} (${Math.round(data.percentage_increase)}%)`
            };
            setStatsData(prevStats => [newUserStat, ...prevStats]);
        })
        .catch(error => console.error('Error:', error));
    }

    const fetchTotalSale = (startDate=defaultStartDate, endDate=defaultEndDate) => {
        axios.post('/api/dashboard/calculate-total-revenue', {
            start_date: startDate,
            end_date: endDate
        })
        .then(response => {
            const data = response.data;
            const arrow = data.is_increase ? '↗︎' : '↙';
            const totalSaleStat = {
                title : "Total Revenue", 
                value : toVND(data.value), 
                icon : <CreditCardIcon className='w-8 h-8'/>, 
                description : `${arrow} ${toVND(data.difference)} (${Math.round(data.percentage_increase)}%)`
            };
            setStatsData(prevStats => [totalSaleStat, ...prevStats]);
        })
        .catch(error => console.error('Error:', error));
    }

    const fetchTotalTrip = (startDate=defaultStartDate, endDate=defaultEndDate) => {
        axios.post('/api/dashboard/calculate-total-trips', {
            start_date: startDate,
            end_date: endDate
        })
        .then(response => {
            const data = response.data;
            const arrow = data.is_increase ? '↗︎' : '↙';
            const totalTripStat = {
                title : "Total Trips", 
                value : data.value, 
                icon : <CircleStackIcon className='w-8 h-8'/>, 
                description : `${arrow} ${data.difference} (${Math.round(data.percentage_increase)}%)`
            };
            setStatsData(prevStats => [totalTripStat, ...prevStats]);
        })
        .catch(error => console.error('Error:', error));
    }

    useEffect(() => {
        fetchTotalTrip();
        fetchTotalSale();    
        fetchNewUsers();  
    }, []);

    const dispatch = useDispatch()
 

    const updateDashboardPeriod = (newRange) => {
        // Dashboard range changed, write code to refresh your values
        dispatch(showNotification({message : `Period updated to ${newRange.startDate} to ${newRange.endDate}`, status : 1}))
        setStatsData([{title : "Active Users", value : Math.round(Math.random() * 1000), icon : <UsersIcon className='w-8 h-8'/>, description : `↙ ${Math.round(Math.random() * 100)} (${Math.round(Math.random() * 10)}%)`}]);
        fetchTotalTrip(newRange.startDate, newRange.endDate); 
        fetchTotalSale(newRange.startDate, newRange.endDate);
        fetchNewUsers(newRange.startDate, newRange.endDate);
    }

    return(
        <>
        {/** ---------------------- Select Period Content ------------------------- */}
            <DashboardTopBar updateDashboardPeriod={updateDashboardPeriod}/>
        
        {/** ---------------------- Different stats content 1 ------------------------- */}
            <div className="grid lg:grid-cols-4 mt-2 md:grid-cols-2 grid-cols-1 gap-6">
                {
                    statsData.map((d, k) => {
                        return (
                            <DashboardStats key={k} {...d} colorIndex={k}/>
                        )
                    })
                }
            </div>



        {/** ---------------------- Different charts ------------------------- */}
            <div className="grid lg:grid-cols-2 mt-4 grid-cols-1 gap-6">
                <BarChartCustom />
                <LineChart />
            </div>
            
        {/** ---------------------- Different stats content 2 ------------------------- */}
        
            {/* <div className="grid lg:grid-cols-2 mt-10 grid-cols-1 gap-6">
                <AmountStats />
                <PageStats />
            </div> */}

        {/** ---------------------- User source channels table  ------------------------- */}
        
            {/* <div className="grid lg:grid-cols-2 mt-4 grid-cols-1 gap-6">
                <UserChannels />
                <DoughnutChart />
            </div> */}
        </>
    )
}

export default Dashboard