import { useState } from "react"
import { useDispatch } from "react-redux"
import InputText from '../../../components/Input/InputText'
import ErrorText from '../../../components/Typography/ErrorText'
import { showNotification } from "../../common/headerSlice"
import { addNewLead } from "../leadSlice"
import { createAsyncThunk } from '@reduxjs/toolkit'
import axios from "axios"
import { getTrips } from "../leadSlice"

export const BookBizopsTrip = createAsyncThunk(
    'leads/addNewLead',
    async (newLeadObj) => {
      const response = await axios.post('/api/create-bizops-trip', newLeadObj);
      const data = response.data;
        localStorage.setItem('current_trip_id', data.trip_id)
      return data
    }
  )

const vehicle = ['Urep Bike', 'Urep Car 4', 'Urep Car 7', 'Urep Car Premium'];

// Use state to keep track of the selected option


const INITIAL_BOOKING_OBJ = {
    user_phone: "",
    user_name: "",
    vehicle: "",
    departure_point: {
        lat: "",
        lng: "", // TODO: update by google map
    },
    departure_name: "",
    destination_point: {
        lat: "",
        lng: "",
    },
    destination_name: "",
    fare: ""
}

function AddLeadModalBody({closeModal}){
    const dispatch = useDispatch()
    const [loading, setLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState("")
    const [bookingObj, setBookingObj] = useState(INITIAL_BOOKING_OBJ)
    const [selectedOption, setSelectedOption] = useState(vehicle[0]);

    const handleSelectChange = (e) => {
        setSelectedOption(e.target.value);
    };
    
    const bookNewBizopsTrip = () => {
        let vehicle_id = vehicle.findIndex(item => item === selectedOption) + 1;
        localStorage.setItem('vehicle_id', vehicle_id)
        if(bookingObj.user_name.trim() === "")return setErrorMessage("Name of User is required!") 
        else if(bookingObj.user_phone.trim() === "")return setErrorMessage("Phone number is required!") 
        else if(bookingObj.departure_name.trim() === "")return setErrorMessage("Departure location is required!") 
        else if(bookingObj.destination_name.trim() === "")return setErrorMessage("Destination location is required!")
        else if(bookingObj.fare.trim() === "")return setErrorMessage("Fare is required!") 
        else{
            let newBookingObj = {
                "user_name": bookingObj.user_name,
                "user_phone": bookingObj.user_phone,
                "vehicle": vehicle_id,
                "departure_point": {
                    "lat": 10.76757020798936,
                    "lng": 106.66339882449851,
                },
                "departure_name": bookingObj.departure_name,
                "destination_point": {
                    "lat": 10.786223457215337,
                    "lng": 106.69830073799115,
                },
                "destination_name": bookingObj.destination_name,
                "fare": parseInt(bookingObj.fare)
            }
            dispatch(BookBizopsTrip(newBookingObj))
            dispatch(showNotification({message : "Booking New Trip Done~", status : 1}))
            closeModal()
        }
    }

    // const saveNewLead = () => {
    //     if(leadObj.first_name.trim() === "")return setErrorMessage("First Name is required!")
    //     else if(leadObj.email.trim() === "")return setErrorMessage("Email id is required!")
    //     else{
    //         let newLeadObj = {
    //             "id": 7,
    //             "email": leadObj.email,
    //             "first_name": leadObj.first_name,
    //             "last_name": leadObj.last_name,
    //             "avatar": "https://reqres.in/img/faces/1-image.jpg"
    //         }
    //         dispatch(addNewLead({newLeadObj}))
    //         dispatch(showNotification({message : "New Lead Added!", status : 1}))
    //         closeModal()
    //     }
    // }

    const updateFormValue = ({updateType, value}) => {
        setErrorMessage("")
        setBookingObj({...bookingObj, [updateType] : value})
    }

    return(
        <>

            <InputText type="text" defaultValue={bookingObj.user_name} updateType="user_name" containerStyle="mt-4" labelTitle="User Name" updateFormValue={updateFormValue}/>

            <InputText type="phone" defaultValue={bookingObj.user_phone} updateType="user_phone" containerStyle="mt-4" labelTitle="Phone Number" updateFormValue={updateFormValue}/>

            <div className={`form-control w-full mt-4`}>
            <label className="label">
                <span className={"label-text text-base-content Vehicle"}>Vehicle</span>
            </label>
                <select 
                    value={selectedOption} 
                    onChange={handleSelectChange}
                    className=""
                    style={{ width: '100%', height: '40px', fontSize: '16px', borderRadius: '7px' }}
                >
                    {vehicle.map((option, index) => (
                        <option key={index} value={option}>
                            {option}
                        </option>
                    ))}
                </select>
            </div>

            <InputText type="text" defaultValue={bookingObj.departure_name} updateType="departure_name" containerStyle="mt-4" labelTitle="Departure Address" updateFormValue={updateFormValue}/>

            <InputText type="text" defaultValue={bookingObj.departure_name} updateType="destination_name" containerStyle="mt-4" labelTitle="Destination Address" updateFormValue={updateFormValue}/> 

            <InputText type="text" defaultValue={bookingObj.fare} updateType="fare" containerStyle="mt-4" labelTitle="Fare" updateFormValue={updateFormValue}/> 

            {/* <InputText type="email" defaultValue={bookingObj.email} updateType="email" containerStyle="mt-4" labelTitle="Email Id" updateFormValue={updateFormValue}/> */}


            <ErrorText styleClass="mt-16">{errorMessage}</ErrorText>
            <div className="modal-action">
                <button  className="btn btn-ghost" onClick={() => closeModal()}>Cancel</button>
                <button  className="btn btn-primary px-6" onClick={() => bookNewBizopsTrip()}>Save</button>
            </div>
        </>
    )
}

export default AddLeadModalBody