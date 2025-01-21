export default function Loading() {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="relative w-24 h-24">
          <div className="absolute top-0 left-0 right-0 bottom-0">
            <div className="w-full h-full border-4 border-gray-200 rounded-full"></div>
            <div className="w-full h-full border-4 border-black rounded-full border-t-transparent animate-spin"></div>
          </div>
        </div>
      </div>
    );
  }