package main

import (
    "log"
    "math/rand"
    "sync"
    "sync/atomic"
    "time"
)

type (
    semaphore chan struct{}
)

type (
    readerWriter struct {
        name string
        write sync.WaitGroup
        readerControl semaphore
        shutdown chan struct{}
        reportShutdown sync.WaitGroup
        maxReads int
        maxReaders int
        currentReads int32
    }
)

func init()  {
    rand.Seed(time.Now().Unix())
}

func main()  {
    log.Println("Starting Process")

    // Create a new readerWriter with a max of 3 reads at a time
    // and a total of 6 reader goroutines.
    first := start("First", 3, 6)

    // Create a new readerWriter with a max of 1 reads at a time
    // and a total of 1 reader goroutines.
    second := start("Second", 2, 2)

    // Let the program run for 2 seconds
    time.Sleep(2 * time.Second)

    // Shutdown both of the readerWriter processes
    shutdown(first, second)

    log.Println("Process Ended")
    return
}

func start(name string, maxReads, maxReaders int) *readerWriter {
    rw := readerWriter{
        name:           name,
        readerControl:  make(semaphore, maxReads),
        shutdown:       make(chan struct{}),
        maxReads:       maxReads,
        maxReaders:     maxReaders,
    }

    rw.reportShutdown.Add(maxReaders)
    for goroutine := 0; goroutine < maxReaders; goroutine++{
        go rw.reader(goroutine)
    }

    rw.reportShutdown.Add(1)
    go rw.writer()

    return &rw
}

// shutdown stops all of the existing readerWriter processes concurrently
func shutdown(readerWriters ...*readerWriter)  {
    // Create a WaitGroup to track the shutdowns
    var waitShutdown sync.WaitGroup
    waitShutdown.Add(len(readerWriters))

    for _, readerWriter := range readerWriters{
        go readerWriter.stop(&waitShutdown)
    }

    // wait for all the goroutines to report they are done
    waitShutdown.Wait()
}

// stop signals to all goroutines to shutdown and reports back
// when that is complete
func (rw *readerWriter) stop(waitShutdown *sync.WaitGroup)  {
    // Schedule the call to Done for once the method returns.
    defer waitShutdown.Done()

    log.Printf("%s\t: #####> Stop", rw.name)

    // Close the channel which will causes all the goroutines waiting on
    // this channel to receive the notification to shutdown
    close(rw.shutdown)

    rw.reportShutdown.Wait()

    log.Printf("%s\t: #####> Stopped", rw.name)
}

// reader is a goroutine that listens on the shutdown channel and
// performs reads until the channel is signaled.
func (rw *readerWriter) reader(reader int)  {
    defer rw.reportShutdown.Done()

    for {
        select {
        case <-rw.shutdown:
            log.Printf("%s\t: #> Reader Shutdown", rw.name)
            return
        default:
            rw.performRead(reader)
        }
    }
}

// performRead performs the actual reading work.
func (rw *readerWriter) performRead(reader int) {
    // Get a read lock for this critical section.
    rw.ReadLock(reader)

    // Safely increment the current reads counter
    count := atomic.AddInt32(&rw.currentReads, 1)

    // Simulate some reading work
    log.Printf("%s\t: [%d] Start\t- [%d] Reads\n", rw.name, reader, count)
    time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

    // Safely decrement the current reads counter
    count = atomic.AddInt32(&rw.currentReads, -1)
    log.Printf("%s\t: [%d] Finish\t- [%d] Reads\n", rw.name, reader, count)

    // Release the read lock for this critical section.
    rw.ReadUnlock(reader)
}

func (rw *readerWriter) writer()  {
    defer rw.reportShutdown.Done()

    for {
        select {
        case <-rw.shutdown:
            log.Printf("%s\t: #> Writer Shutdown", rw.name)
            return
        default:
            rw.performWrite()
        }
    }
}

// performWrite performs the actual write work.
func (rw *readerWriter) performWrite()  {
    // Wait a random number of milliseconds before we write again.
    time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

    log.Printf("%s\t: *****> Writing Pending\n", rw.name)

    // Get a write lock for this critical section.
    rw.WriteLock()

    // Simulate some writing work.
    log.Printf("%s\t: *****> Writing Start", rw.name)
    time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
    log.Printf("%s\t: *****> Writing Finish", rw.name)

    // Release the write lock for this critical section.
    rw.WriteUnlock()
}

func (rw *readerWriter) ReadLock(reader int)  {
    rw.write.Wait()

    rw.readerControl.Acquire(1)
}

func (rw *readerWriter) ReadUnlock(reader int)  {
    rw.readerControl.Release(1)
}

// WriteLock blocks all reading so the write can happen safely.
func (rw *readerWriter) WriteLock() {
    // Add 1 to the waitGroup so reads will block
    rw.write.Add(1)

    // Acquire all the buffers from the semaphore channel.
    rw.readerControl.Acquire(rw.maxReads)
}

// WriteUnlock releases the write lock and allows reads to occur.
func (rw *readerWriter) WriteUnlock() {
    // Release all the buffers back into the semaphore channel.
    rw.readerControl.Release(rw.maxReads)

    // Release the write lock.
    rw.write.Done()
}

func (s semaphore) Acquire(buffers int)  {
    var e struct{}

    for buffer := 0; buffer < buffers; buffer++{
        s <- e
    }
}

// Release returns the specified number of buffers back into the semaphore channel.
func (s semaphore) Release(buffers int) {
    // Read the data from the channel to release buffers.
    for buffer := 0; buffer < buffers; buffer++ {
        <-s
    }
}