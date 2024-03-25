import {
  createSignal,
  createResource,
  // Match,
  Show,
  For,
  // Switch,
  JSX,
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
  const response = await fetch(`http://${serveraddress}/thumbnail/${device}`, {
    cache: "no-cache",
  });
  return response.blob();
}

async function fetchMediaFileList(): Promise<string[]> {
  const response = await fetch(`http://${serveraddress}/mediafilelist`, {
    cache: "no-cache",
  });

  return response.json();
}

async function createNewThumbnail(device: string): Promise<void> {
  await fetch(`http://${serveraddress}/createthumbnail/${device}`, {
    cache: "no-cache",
  });
}

async function startRecording(device: string): Promise<response> {
  const response = await fetch(
    `http://${serveraddress}/startcapture/${device}`,
    {
      cache: "no-cache",
    },
  );
  const responseJson: response = await response.json();
  if (responseJson.success) {
    console.log("Recording started");
  } else {
    window.alert(responseJson.message);
  }
  return responseJson;
}

async function stopRecording(): Promise<response> {
  const response = await fetch(`http://${serveraddress}/stopcapture`, {
    cache: "no-cache",
  });
  return response.json();
}

async function fetchActiveRecordingStatus(): Promise<response> {
  const response = await fetch(`http://${serveraddress}/activerecording`, {
    cache: "no-cache",
  });
  return response.json();
}

function Button(props: {
  onClick: () => void;
  text: string;
  color: string;
}): JSX.Element {
  const buttonClass = `border-1 ${props.color} p-1`;
  return (
    <button class={buttonClass} onClick={props.onClick}>
      {props.text}
    </button>
  );
}

function App() {
  const [videoDevices] = createResource(fetchVideoDevices);
  const [videoDevice, setVideoDevice] = createSignal("video0");
  const [videoThumbnail, { refetch: refetchVideoThumbnail }] = createResource(
    videoDevice,
    fetchVideoThumbnail,
  );
  const [activeRecordingStatus, { refetch: refetchActiveRecordingStatus }] =
    createResource(fetchActiveRecordingStatus);
  const [mediaFileList, { refetch: refetchMediaFileList }] =
    createResource(fetchMediaFileList);

  (async () => {
    setInterval(() => {
      refetchActiveRecordingStatus();
      refetchMediaFileList();
      if (activeRecordingStatus()?.success) {
        refetchVideoThumbnail();
      }
    }, 3000);
  })();

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
          <Button
            color="bg-slate-600"
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
            text="Refresh Thumbnail"
          />
        </div>
        <div>
          <Show when={activeRecordingStatus()}>
            <div>
              {activeRecordingStatus()?.success ? (
                <Button
                  color="bg-amber-600"
                  onClick={stopRecording}
                  text="Stop Recording"
                />
              ) : (
                <Button
                  color="bg-red-600"
                  onClick={() => startRecording(videoDevice())}
                  text="Start Recording"
                />
              )}
            </div>
          </Show>
        </div>
      </Show>
    </>
  );
}

export default App;
