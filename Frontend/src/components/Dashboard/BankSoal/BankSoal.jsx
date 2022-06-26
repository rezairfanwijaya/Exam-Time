import React, { useEffect, useState } from "react";
import Card from "./CardBank/Card";
import axios from "axios";

const BankSoal = () => {
  const [mapel, setMapel] = useState()
  const [image, setImage] = useState()
  const [name, setName] = useState()

  const fetchMapel = async () => {
    try {
      const res = await axios.get('http://localhost:8080/api/guru/dashboard', 
        {
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
          }
        }
      )
      
     
      const mataPelajaran = res.data.data
      console.log("Berhasil fetch data", mataPelajaran)
      //sethere
      
      setName(res.data.data.nama)
      setImage(res.data.data.avatar)
      setMapel(mataPelajaran)

    } catch (error) {
      console.log("Gagal fetch data mapel", error)
    }
  }
  
  useEffect(() => {
    fetchMapel()
  }, [])


  return (
    <>
      <div className="flex justify-start">
        <div className="grid grid-cols-3 mr-5">
          {mapel && mapel.map((mapel) => {
            return <Card mapel={mapel.mata_pelajaran} id_mapel={mapel.id_mata_pelajaran} />
          })}
        </div>
      </div>
    </>
  );
};

export default BankSoal;
