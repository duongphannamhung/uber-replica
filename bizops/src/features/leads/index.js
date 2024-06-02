import moment from "moment"
import { useRef, useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import TitleCard from "../../components/Cards/TitleCard"
import { openModal } from "../common/modalSlice"
import { deleteLead, getTrips } from "./leadSlice"
import { CONFIRMATION_MODAL_CLOSE_TYPES, MODAL_BODY_TYPES } from '../../utils/globalConstantUtil'
import TrashIcon from '@heroicons/react/24/outline/TrashIcon'
import { showNotification } from '../common/headerSlice'
import { useState } from "react"
import Pagination from 'react-js-pagination';
import axios from "axios"
import { createAsyncThunk } from '@reduxjs/toolkit'
import ArrowPathIcon  from '@heroicons/react/24/outline/ArrowPathIcon'



const findDriver = async (intervalId) => {
    await axios.get(`http://localhost:6969/api/bizops/find-driver?trip_id=${localStorage.getItem('current_trip_id')}&vehicle_type=${localStorage.getItem('vehicle_id')}`)
        .then(async (response) => {
            if (response.data.find_done) {
                clearInterval(intervalId.current);
                localStorage.removeItem('current_trip_id')
                localStorage.removeItem('vehicle_id')
            }
        })
        .catch((error) => {
            console.error(error)
            alert(error.response.data.message)
        })
}

const TopSideButtons = () => {
    const dispatch = useDispatch()

    const bookBizopsTrip = () => {
        dispatch(openModal({title : "Book Bizops Trip", bodyType : MODAL_BODY_TYPES.LEAD_ADD_NEW}))
    }

    return(
        <div className="inline-block float-right">
            <button 
                className="btn btn-ghost btn-sm normal-case"
                onClick={() => dispatch(getTrips(1))}
            >
                <ArrowPathIcon className="w-4 mr-2 mt=2"/>Refresh
            </button>
            <button className="btn px-6 btn-sm normal-case btn-primary" onClick={() => bookBizopsTrip() }>Book trip</button>
        </div>
    )
}

function Leads(){
    // const [totalItemsCount, setTotalItemsCount] = useState(0); // Add this line to store the total number of items
    const {leads } = useSelector(state => state.lead)
    const dispatch = useDispatch()
    const [page, setPage] = useState(1);
    const intervalId = useRef();

    useEffect(() => {
        dispatch(getTrips(page))
    }, [page])

    useEffect(() => {
        if(localStorage.getItem('current_trip_id')) {
            intervalId.current = setInterval(() => findDriver(intervalId), 10000);
    }
        return () => clearInterval(intervalId.current)
    }, [])

    // const fetchTrips = async () => {
        // const response = await getTrips(page);
        // setTotalItemsCount(response.total);
    // };

    // useEffect(() => {
    //     // Modify this function to also set the total number of items
    //     // This assumes getTrips returns an object with properties 'items' and 'totalItemsCount'
    //     // fetchTrips();
    //     getTrips(page);
    //     // console.log("totalItemsCount", totalItemsCount)
    // }, [page]);


    const getDummyStatus = (index) => {
        if(index % 5 === 0)return <div className="badge">Not Interested</div>
        else if(index % 5 === 1)return <div className="badge badge-primary">In Progress</div>
        else if(index % 5 === 2)return <div className="badge badge-secondary">Sold</div>
        else if(index % 5 === 3)return <div className="badge badge-accent">Need Followup</div>
        else return <div className="badge badge-ghost">Open</div>
    }

    const getVehicleType = (index) => {
        if(index === 1)return <div className="badge badge-secondary">UrepBike</div>
        else if(index === 2)return <div className="badge badge-primary">UrepCar</div>
        else if(index === 3)return <div className="badge badge-accent">UrepCar7</div>
        else return <div className="badge badge-info">UrepPlus</div>
    }

    const getIsStarted = (is_started) => {
        if(is_started) return <div className="badge badge-success">Yes</div>
        else return <div className="badge badge-ghost">No</div>
    }

    const _deleteTrip = createAsyncThunk('/leads/delete', async (index) => {
    const response = await axios.get(`/api/delete-trip/${index}`, {})
        return response;
    })

    const deleteTrip = (index) => {
        dispatch(openModal({title : "Confirmation", bodyType : MODAL_BODY_TYPES.CONFIRMATION, 
        extraObject : { message : `Are you sure you want to delete Trip id ${index}?`, type : CONFIRMATION_MODAL_CLOSE_TYPES.TRIP_DELETE, index}}))
        dispatch(_deleteTrip(index));
    }

    return(
        <>
            <TitleCard title="Current Trips" topMargin="mt-2" TopSideButtons={<TopSideButtons />}>
                {/* Leads List in table format loaded from slice after api call */}
            <div className="overflow-x-auto w-full">
                <table className="table w-full">
                    <thead>
                    <tr>
                        <th>Trip Id</th>
                        <th>Driver Id</th>
                        <th>Service Type</th>
                        <th>Is Started</th>
                        <th>Departure Place</th>
                        <th>Destination Place</th>
                        <th>Fare</th>
                        <th>Time Created</th>
                        <th></th>
                    </tr>
                    </thead>
                    <tbody>
                        {
                            leads.map((l, k) => {
                                return(
                                    <tr key={k}>
                                    <td style={{ textAlign: 'center' }}>{l.id}</td>
                                    <td style={{ textAlign: 'center' }}>{l.driver_id.Int32}</td>
                                    <td style={{ textAlign: 'center' }}>{getVehicleType(l.service_type)}</td>
                                    <td style={{ textAlign: 'center' }}>{getIsStarted(l.is_started)}</td>
                                    <td style={{ textAlign: 'center' }}>{l.departure_name}</td>
                                    <td style={{ textAlign: 'center' }}>{l.destination_name}</td>
                                    <td style={{ textAlign: 'center' }}>{l.fare.Int32}</td>
                                    <td style={{ textAlign: 'center' }}>{l.created_at}</td>
                                    <td><button className="btn btn-square btn-ghost" onClick={() => deleteTrip(l.id)}><TrashIcon className="w-5"/></button></td>
                                    </tr>
                                )
                            })
                        }
                    </tbody>
                </table>
            </div>
            </TitleCard>
            <style>
                {`
                    .pagination li {
                    display: inline-block !important; margin-right: 10px; /* Add space to the right of each page link */
                    }
                    .active-page-link {
                        font-weight: bold !important;
                        color: ; /* Set the text color */
                        border: 1px solid currentColor; /* Use the current text color for the border */
                        padding: 1px;
                      }
                `}
                </style>
                <div style={{ display: 'flex', justifyContent: 'center' }}>
                <Pagination
                    activePage={page}
                    itemsCountPerPage={10}
                    totalItemsCount={5000}
                    pageRangeDisplayed={7}
                    onChange={(pageNumber) => setPage(pageNumber)}
                    activeLinkClass="active-page-link"
                />
                </div>
        </>
    )
}

{/* <td>{moment(new Date()).add(-5*(k+2), 'days').format("DD MMM YY")}</td> 
            <button disabled={page === 1} onClick={() => setPage(page - 1)}>
                &lt;
            </button>

            Add clickable numbers for the offset here
            {[...Array(10).keys()].slice(0, 3).map(i =>
                <button key={i} onClick={() => setPage(i + 1)}>
                    {i + 1}
                </button>
            )}
            <button onClick={() => setPage(page + 1)}>
                &gt;
            </button>

*/}

export default Leads