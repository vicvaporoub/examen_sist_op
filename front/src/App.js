import React, { useEffect, useState } from 'react';

function App() {
  const [alumnos, setAlumnos] = useState(''); // Estado para guardar texto de la API

  useEffect(() => {
    fetch("http://localhost:3000/alumnos")
 
      .then(response => response.text())   // Le dices que esperas texto plano
      .then(data => setAlumnos(data));     // Guardas el resultado
  }, []);

  return (
    <div className="App">
      <h1>Listado de alumnos</h1>
      <pre>{alumnos}</pre> {/* Muestra el texto recibido */}
    </div>
  );
}

export default App;
