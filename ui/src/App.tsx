import {
  createSignal,
  createResource,
  // Match,
  Show,
  For,
  // Switch,
} from "solid-js";
import "./App.css";

type response = {
  success: boolean;
  message: string;
};

// get ip of server
const serveraddress = window.location.hostname;

async function fetchVideoDevices(): Promise<string[]> {
  const response = await fetch(`http://${serveraddress}/videodevices`);
  return response.json();
}

async function fetchVideoThumbnail(device: string): Promise<Blob> {
  // get jpg from /media/thumbnails/{device}
  const response = await fetch(
    `http://${serveraddress}/media/thumbnails/${device}.jpg`,
  );
  return response.blob();
}

async function createNewThumbnail(device: string): Promise<void> {
  await fetch(`http://${serveraddress}/createthumbnail/${device}`);
}

async function fetchActiveRecording(): Promise<response> {
  const response = await fetch(`http://${serveraddress}/activerecording`);
  return response.json();
}

function App() {
  const [videoDevices] = createResource(fetchVideoDevices);
  const [videoDevice, setVideoDevice] = createSignal("video0");
  const [videoThumbnail, { refetch: refetchVideoThumbnail }] = createResource(
    videoDevice,
    fetchVideoThumbnail,
  );
  const [activeRecording] = createResource(fetchActiveRecording);

  return (
    <>
      <Show when={videoDevices()}>
        <div>
          <div>
            <label for="vidselect">Select video device:</label>
            <select
              onChange={(e) => {
                setVideoDevice(e.currentTarget.value);
              }}
              id="vidselect"
            >
              <For each={videoDevices()}>
                {(device) =>
                  device === "video0" ? (
                    <option selected={true} value={device}>
                      {device}
                    </option>
                  ) : (
                    <option value={device}>{device}</option>
                  )
                }
              </For>
            </select>
          </div>
        </div>
        <div>
          <Show when={videoThumbnail()}>
            <img
              src={URL.createObjectURL(videoThumbnail() as Blob)}
              alt="video thumbnail"
            />
          </Show>
          <button
            // class="p-1"
            onClick={async () => {
              createNewThumbnail(videoDevice());
              (() => {
                setTimeout(refetchVideoThumbnail, 1000);
                refetchVideoThumbnail;
                setTimeout(refetchVideoThumbnail, 2000);
                refetchVideoThumbnail;
                setTimeout(refetchVideoThumbnail, 3000);
                refetchVideoThumbnail;
                setTimeout(refetchVideoThumbnail, 4000);
                refetchVideoThumbnail;
                setTimeout(refetchVideoThumbnail, 5000);
                refetchVideoThumbnail;
              })();
            }}
          >
            Refresh Thumbnail
          </button>
        </div>
        <div>
          <Show when={activeRecording()}>
            <div>
              {activeRecording()?.success
                ? "Recording is active"
                : "Recording is not active"}
            </div>
          </Show>
        </div>
      </Show>
    </>
  );
}

export default App;
