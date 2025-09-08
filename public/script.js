let juegos = [];

const formJuego = document.getElementById("form-juego");
const juegoTitulo = document.getElementById("juego-titulo");

const formMapa = document.getElementById("form-mapa");
const mapaTitulo = document.getElementById("mapa-titulo");
const selectJuegoMapa = document.getElementById("select-juego-mapa");

const selectJuegoRandom = document.getElementById("select-juego-random");
const listaMapas = document.getElementById("lista-mapas");

const btnMapaAleatorio = document.getElementById("btn-mapa-aleatorio");
const resultadoAleatorio = document.getElementById("resultado-aleatorio");
const btnReiniciarBaneos = document.getElementById("btn-reiniciar-baneos");

function actualizarSelects() {
  [selectJuegoMapa, selectJuegoRandom].forEach(select => {

    const selectedValue = select.value;

    select.innerHTML = '<option value="">-- Selecciona un juego --</option>';

    juegos.forEach(juego => {
      const option = document.createElement("option");
      option.value = juego.id;
      option.textContent = juego.titulo;
      select.appendChild(option);
    });

    if (selectedValue) select.value = selectedValue;
  });
}

function mostrarMapas() {
  listaMapas.innerHTML = "";
  const juegoId = selectJuegoRandom.value;
  if (!juegoId) return;

  const juego = juegos.find(j => j.id == juegoId);
  if (!juego) return;

  juego.mapas.forEach(mapa => {
    const li = document.createElement("li");
    const checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.checked = mapa.baneado;
    checkbox.addEventListener("change", () => {
      mapa.baneado = checkbox.checked;
    });

    li.appendChild(checkbox);
    li.append(` ${mapa.titulo}`);
    listaMapas.appendChild(li);
  });
}

function generarId() {
  return Date.now() + Math.floor(Math.random() * 1000);
}

formJuego.addEventListener("submit", e => {
  e.preventDefault();
  const titulo = juegoTitulo.value.trim();
  if (!titulo) return;

  const nuevoJuego = {
    id: generarId(),
    titulo,
    mapas: []
  };

  juegos.push(nuevoJuego);
  juegoTitulo.value = "";
  actualizarSelects();
});

formMapa.addEventListener("submit", e => {
  e.preventDefault();
  const titulo = mapaTitulo.value.trim();
  const juegoId = selectJuegoMapa.value;
  if (!titulo || !juegoId) return;

  const juego = juegos.find(j => j.id == juegoId);
  if (!juego) return;

  const nuevoMapa = {
    id: generarId(),
    titulo,
    baneado: false
  };

  juego.mapas.push(nuevoMapa);
  mapaTitulo.value = "";

  if (selectJuegoRandom.value == juegoId) {
    mostrarMapas();
  }
});

selectJuegoRandom.addEventListener("change", mostrarMapas);

btnReiniciarBaneos.addEventListener("click", () => {
  const juegoId = selectJuegoRandom.value;
  if (!juegoId) return;
  const juego = juegos.find(j => j.id == juegoId);
  if (!juego) return;
  juego.mapas.forEach(mapa => mapa.baneado = false);
  mostrarMapas();
});

btnMapaAleatorio.addEventListener("click", () => {
  const juegoId = selectJuegoRandom.value;
  if (!juegoId) {
    resultadoAleatorio.textContent = "Selecciona un juego primero.";
    return;
  }
  const juego = juegos.find(j => j.id == juegoId);
  if (!juego) return;

  const mapasDisponibles = juego.mapas.filter(m => !m.baneado);
  if (mapasDisponibles.length === 0) {
    resultadoAleatorio.textContent = "No hay mapas disponibles (todos baneados).";
    return;
  }

  const aleatorio = mapasDisponibles[Math.floor(Math.random() * mapasDisponibles.length)];
  resultadoAleatorio.textContent = `Mapa seleccionado: ${aleatorio.titulo}`;
});
